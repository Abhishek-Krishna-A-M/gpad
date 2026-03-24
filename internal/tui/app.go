// Package tui implements the full-screen gpad TUI.
// Zero external dependencies beyond golang.org/x/term.
//
// File layout:
//   app.go       — App struct, Run(), input loop, mode dispatch
//   draw.go      — coordinate-based rendering engine (no \n in draw path)
//   tree.go      — treeNode, buildTree, lf-style cd navigation
//   navigate.go  — cursor movement, scroll, Enter/Space/h/l
//   actions.go   — open, delete, rename, copy, paste, pin, sync
//   panels.go    — search, tags, graph, daily, links overlay panels
//   preview.go   — markdown rendering for the right-hand panel
//   keybinds.go  — load ~/.gpad/keybinds.json, action dispatch
//   ansi.go      — ANSI constants, mc(), vlen(), truncANSI(), pad()
package tui

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/config"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/gitrepo"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/templates"
)

// ── Modes ─────────────────────────────────────────────────────────────────────

type mode int

const (
	modeNormal      mode = iota
	modeCommand          // : command bar active
	modeFilter           // / tree filter active
	modePanelFilter      // / inside a panel
	modeHelp             // ? help overlay
	modeConfirm          // y/N confirm dialog
)

// ── App struct ────────────────────────────────────────────────────────────────

type App struct {
	// terminal
	width    int
	height   int
	fd       int
	oldState *term.State

	// vault
	notesRoot string
	cwd       string        // current working directory (lf-style cd)
	cwdStack  []string      // stack for cd back (- key)

	// tree state
	treeRoot   *treeNode
	flat       []*treeNode  // flattened visible nodes
	cursor     int
	treeScroll int
	treeWidth  int

	// filter (tree filter mode)
	filterStr string

	// command bar
	cmdBuf  string
	cmdHist []string
	cmdIdx  int

	// confirm dialog
	confirmMsg    string
	confirmAction func() error

	// preview cache — keyed by abs path, invalidated on cursor move
	previewCache string
	previewLines []string

	// panel overlay (search/tags/graph/daily/links)
	panel *PanelState

	// status
	mode      mode
	statusMsg string
	gitStatus string

	// yank buffer
	yankBuf string

	// keybinds
	keyMap KeyMap
}

// ── Entry point ───────────────────────────────────────────────────────────────

func Run() error {
	fd := int(os.Stdin.Fd())
	state, err := makeRaw(fd)
	if err != nil {
		return fmt.Errorf("raw terminal: %w", err)
	}

	w, h, err := term.GetSize(fd)
	if err != nil {
		w, h = 80, 24
	}

	storage.EnsureDirs()
	templates.EnsureDefaults()

	a := &App{
		fd:        fd,
		oldState:  state,
		width:     w,
		height:    h,
		notesRoot: storage.NotesDir(),
		cwd:       storage.NotesDir(),
		treeWidth: 28,
		mode:      modeNormal,
		keyMap:    LoadKeyMap(),
	}

	if err := a.buildTree(); err != nil {
		restoreTerminal(fd, state)
		return err
	}

	a.gitStatus = a.checkGit()
	SaveDefaultKeybinds()

	fmt.Print(aHideCur + aClearScr)
	defer func() {
		fmt.Print(aShowCur + aClearScr + aHome)
		restoreTerminal(fd, state)
	}()

	a.draw()
	return a.loop()
}

// makeRaw switches stdin to raw mode and returns the old state.
func makeRaw(fd int) (*term.State, error) {
	return term.MakeRaw(fd)
}

// restoreTerminal restores stdin to its original state.
func restoreTerminal(fd int, state *term.State) {
	term.Restore(fd, state)
}

// ── Input loop ────────────────────────────────────────────────────────────────

func (a *App) loop() error {
	buf := make([]byte, 32)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			return nil
		}
		b := buf[:n]

		// check terminal resize
		if w, h, e := term.GetSize(a.fd); e == nil {
			if w != a.width || h != a.height {
				a.width, a.height = w, h
				a.treeWidth = a.width / 3
				if a.treeWidth < 18 {
					a.treeWidth = 18
				}
				if a.treeWidth > 42 {
					a.treeWidth = 42
				}
				os.Stdout.WriteString(aClearScr + aHome)
			}
		}

		quit := a.dispatch(b)
		if quit || a.statusMsg == "__quit__" {
			return nil
		}
		a.draw()
	}
}

// ── Dispatch ──────────────────────────────────────────────────────────────────

// dispatch routes raw input bytes to the correct handler based on mode.
func (a *App) dispatch(b []byte) (quit bool) {
	// modes that have their own full input handling
	switch a.mode {
	case modeCommand:
		return a.handleCommandInput(b)
	case modeHelp:
		a.mode = modeNormal
		return false
	case modeConfirm:
		return a.handleConfirmInput(b)
	case modeFilter:
		return a.handleFilterInput(b)
	case modePanelFilter:
		return a.handlePanelFilterInput(b)
	}

	// panel normal mode
	if a.panel != nil {
		return a.handlePanelInput(b)
	}

	// normal mode — look up action in keymap
	key := rawKeyToString(b)
	action, ok := a.keyMap[key]
	if !ok {
		return false
	}
	return a.execAction(action)
}

// execAction executes a named action. Returns true to quit.
func (a *App) execAction(act Action) bool {
	switch act {
	case ActionQuit:
		return true
	case ActionMoveUp:
		a.moveUp()
	case ActionMoveDown:
		a.moveDown()
	case ActionMoveTop:
		a.moveTop()
	case ActionMoveBottom:
		a.moveBottom()
	case ActionEnter:
		a.handleEnter()
	case ActionExpand:
		a.handleExpand()
	case ActionLeft:
		a.collapseOrUp()
	case ActionUp:
		a.cdUp()
	case ActionOpen:
		a.openSelected()
	case ActionView:
		a.viewSelected()
	case ActionNew:
		a.newNote()
	case ActionDelete:
		a.deleteSelected()
	case ActionRename:
		a.renameSelected()
	case ActionYank:
		a.yankSelected()
	case ActionPaste:
		a.pasteYanked()
	case ActionPin:
		a.togglePin()
	case ActionSync:
		a.runSync()
	case ActionToday:
		a.openToday()
	case ActionCommandMode:
		a.mode = modeCommand
		a.cmdBuf = ""
		a.statusMsg = ""
	case ActionFilterMode:
		a.mode = modeFilter
		a.filterStr = ""
	case ActionHelp:
		a.mode = modeHelp
	case ActionRedraw:
		prevW, prevH = 0, 0 // force full clear on next draw
		fmt.Print(aClearScr)
	case ActionPanelSearch:
		a.openPanel(panelSearch)
		a.mode = modePanelFilter
		a.loadSearchResults("")
	case ActionPanelTags:
		a.openPanel(panelTags)
	case ActionPanelGraph:
		a.openPanel(panelGraph)
	case ActionPanelDaily:
		a.openPanel(panelDaily)
	case ActionPanelLinks:
		a.openPanel(panelLinks)
	}
	return false
}

// ── Command mode input ────────────────────────────────────────────────────────

func (a *App) handleCommandInput(b []byte) bool {
	ch := rune(b[0])

	// arrow keys in command mode = history
	if b[0] == 0x1b && len(b) >= 3 && b[1] == '[' {
		switch b[2] {
		case 'A': // up — prev history
			if a.cmdIdx < len(a.cmdHist)-1 {
				a.cmdIdx++
				a.cmdBuf = a.cmdHist[len(a.cmdHist)-1-a.cmdIdx]
			}
		case 'B': // down — next history
			if a.cmdIdx > 0 {
				a.cmdIdx--
				a.cmdBuf = a.cmdHist[len(a.cmdHist)-1-a.cmdIdx]
			} else {
				a.cmdIdx = 0
				a.cmdBuf = ""
			}
		}
		return false
	}

	switch {
	case ch == 0x1b: // Esc
		a.mode = modeNormal
		a.cmdBuf = ""
		a.statusMsg = ""
	case ch == '\r' || ch == '\n':
		a.execCommandStr(a.cmdBuf)
		if a.cmdBuf != "" {
			a.cmdHist = append(a.cmdHist, a.cmdBuf)
		}
		a.cmdIdx = 0
		a.cmdBuf = ""
		a.mode = modeNormal
	case ch == 0x7f || ch == 0x08: // backspace
		if len(a.cmdBuf) == 0 {
			a.mode = modeNormal
		} else {
			_, size := lastRuneSize(a.cmdBuf)
			a.cmdBuf = a.cmdBuf[:len(a.cmdBuf)-size]
		}
	case ch == 0x03: // Ctrl-C
		a.mode = modeNormal
		a.cmdBuf = ""
	case ch >= 0x20:
		a.cmdBuf += string(ch)
	}
	return false
}

// execCommandStr parses and executes a colon-command string.
func (a *App) execCommandStr(cmd string) {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return
	}
	parts := strings.Fields(cmd)
	verb, args := parts[0], parts[1:]

	switch verb {
	case "q", "quit":
		// mark quit — handled by returning true from dispatch on next cycle
		// we set a flag so the loop exits cleanly
		a.statusMsg = "__quit__"
	case "open", "e", "edit":
		if len(args) > 0 {
			abs := storage.AbsPath(args[0])
			a.openAbsPath(abs, args[0])
		} else {
			a.openSelected()
		}
	case "view", "v":
		if len(args) > 0 {
			abs := storage.AbsPath(args[0])
			a.previewCache = ""
			a.previewLines = renderPreview(abs, 9999)
		} else {
			a.viewSelected()
		}
	case "new":
		name, tmpl := "", ""
		for i, arg := range args {
			if arg == "-t" && i+1 < len(args) {
				tmpl = args[i+1]
			} else if !strings.HasPrefix(arg, "-") && name == "" {
				name = arg
			}
		}
		if tmpl != "" {
			_ = tmpl // template support: pass to notes.Create when wired
		}
		a.createNote(name, tmpl)
	case "mkdir":
		if len(args) > 0 {
			a.createDir(args[0])
		}
	case "rm", "delete", "remove":
		if len(args) > 0 {
			a.confirmMsg = "delete " + args[0] + "? [y/N]"
			a.confirmAction = func() error {
				return nil // core.Delete wired via confirmDelete
			}
			_ = a.buildTree()
		} else {
			a.deleteSelected()
		}
	case "mv", "move", "rename":
		if len(args) >= 2 {
			a.runMv(args[0], args[1])
		} else {
			a.setStatus("usage: mv <src> <dest>")
		}
	case "cp", "copy":
		if len(args) >= 2 {
			a.runCp(args[0], args[1])
		} else {
			a.setStatus("usage: cp <src> <dest>")
		}
	case "sync":
		a.runSync()
	case "pin":
		a.togglePin()
	case "links", "backlinks":
		a.openPanel(panelLinks)
	case "tags":
		a.openPanel(panelTags)
	case "graph":
		a.openPanel(panelGraph)
	case "daily":
		a.openPanel(panelDaily)
	case "find", "search":
		q := strings.Join(args, " ")
		a.openPanel(panelSearch)
		a.loadSearchResults(q)
		if q != "" {
			a.mode = modeNormal
		} else {
			a.mode = modePanelFilter
		}
	case "today":
		a.openToday()
	case "help", "?":
		a.mode = modeHelp
	case "keybinds":
		abs := keybindsPath()
		a.openAbsPath(abs, "keybinds.json")
		// reload after editing
		a.keyMap = LoadKeyMap()
	// ── git commands ──────────────────────────────────────────────────
	case "git":
		if len(args) >= 2 && args[0] == "init" {
			url := args[1]
			a.setStatus("initialising git…")
			a.draw()
			if err := gitrepo.Initialize(a.notesRoot, url); err != nil {
				a.setStatus("git init error: " + err.Error())
				return
			}
			cfg, _ := config.Load()
			cfg.GitEnabled = true
			cfg.RepoURL = url
			cfg.AutoPush = true
			_ = config.Save(cfg)
			a.gitStatus = a.checkGit()
			a.setStatus("git connected → " + url)
		} else if len(args) == 1 && args[0] == "status" {
			cfg, _ := config.Load()
			if cfg.GitEnabled {
				a.setStatus("remote: " + cfg.RepoURL + "  autopush: " + boolStr(cfg.AutoPush))
			} else {
				a.setStatus("git not configured — run :git init <url>")
			}
		} else {
			a.setStatus("usage: git init <url>  |  git status")
		}

	// ── config commands ─────────────────────────────────────────────────
	case "config":
		if len(args) < 2 {
			cfg, _ := config.Load()
			a.setStatus(fmt.Sprintf("editor:%s  autopush:%s  git:%s",
				cfg.Editor, boolStr(cfg.AutoPush), boolStr(cfg.GitEnabled)))
			return
		}
		switch args[0] {
		case "editor":
			cfg, _ := config.Load()
			cfg.Editor = args[1]
			_ = config.Save(cfg)
			a.setStatus("editor set to: " + args[1])
		case "autopush":
			cfg, _ := config.Load()
			cfg.AutoPush = args[1] == "on"
			_ = config.Save(cfg)
			a.setStatus("autopush: " + args[1])
		default:
			a.setStatus("usage: config editor <n>  |  config autopush on/off")
		}

	// ── template commands ───────────────────────────────────────────────
	case "template":
		if len(args) == 0 || args[0] == "list" {
			list := templates.List()
			a.setStatus("templates: " + strings.Join(list, ", "))
		} else if args[0] == "new" && len(args) > 1 {
			if err := templates.Save(args[1], defaultTemplateContent(args[1])); err != nil {
				a.setStatus("template error: " + err.Error())
				return
			}
			a.openAbsPath(templates.Path(args[1]), "templates/"+args[1]+".md")
		} else if args[0] == "edit" && len(args) > 1 {
			a.openAbsPath(templates.Path(args[1]), "templates/"+args[1]+".md")
		} else if args[0] == "delete" && len(args) > 1 {
			if err := templates.Delete(args[1]); err != nil {
				a.setStatus("template error: " + err.Error())
				return
			}
			a.setStatus("deleted template: " + args[1])
		} else {
			a.setStatus("usage: template list | new <n> | edit <n> | delete <n>")
		}

	default:
		a.setStatus("unknown: " + verb + " — type ? for help")
	}
}

func boolStr(b bool) string {
	if b { return "on" }
	return "off"
}

func defaultTemplateContent(name string) string {
	return "---\ntitle: {{title}}\ndate: {{date}}\ntags: []\n---\n\n# {{title}}\n\n{{cursor}}\n"
}

// ── Filter mode input ─────────────────────────────────────────────────────────

func (a *App) handleFilterInput(b []byte) bool {
	ch := rune(b[0])
	switch {
	case b[0] == 0x1b:
		a.mode = modeNormal
		a.filterStr = ""
		_ = a.buildTree()
	case ch == '\r' || ch == '\n':
		a.mode = modeNormal
		a.handleEnter()
	case ch == 0x7f || ch == 0x08:
		if len(a.filterStr) == 0 {
			a.mode = modeNormal
			_ = a.buildTree()
		} else {
			_, size := lastRuneSize(a.filterStr)
			a.filterStr = a.filterStr[:len(a.filterStr)-size]
			a.applyTreeFilter()
		}
	case ch >= 0x20:
		a.filterStr += string(ch)
		a.applyTreeFilter()
	}
	return false
}

func (a *App) applyTreeFilter() {
	_ = a.buildTree()
	if a.filterStr == "" {
		return
	}
	pat := strings.ToLower(a.filterStr)
	var filtered []*treeNode
	for _, n := range a.flat {
		if fuzzyMatch(pat, strings.ToLower(n.name)) {
			filtered = append(filtered, n)
		}
	}
	a.flat = filtered
	a.cursor = 0
	a.treeScroll = 0
}

// ── Panel input ───────────────────────────────────────────────────────────────

func (a *App) handlePanelInput(b []byte) bool {
	key := rawKeyToString(b)
	switch key {
	case "j", "down":
		a.panelMoveDown()
	case "k", "up":
		a.panelMoveUp()
	case "g":
		a.panel.cursor = 0
		a.panel.scroll = 0
	case "G":
		a.panel.cursor = len(a.panel.results) - 1
		a.panelScrollIntoView()
	case "enter":
		a.panelEnter()
	case "/":
		a.mode = modePanelFilter
		a.panel.query = ""
		a.panelFilter("")
	case "esc", "q":
		a.closePanel()
		a.mode = modeNormal
	}
	return false
}

func (a *App) handlePanelFilterInput(b []byte) bool {
	ch := rune(b[0])
	switch {
	case b[0] == 0x1b:
		a.mode = modeNormal
		if a.panel != nil {
			a.panel.query = ""
			a.panelFilter("")
		}
	case ch == '\r' || ch == '\n':
		a.mode = modeNormal
		a.panelEnter()
	case ch == 0x7f || ch == 0x08:
		if a.panel != nil && len(a.panel.query) > 0 {
			_, size := lastRuneSize(a.panel.query)
			a.panel.query = a.panel.query[:len(a.panel.query)-size]
			a.panelFilter(a.panel.query)
		} else {
			a.mode = modeNormal
		}
	case ch >= 0x20:
		if a.panel != nil {
			a.panel.query += string(ch)
			a.panelFilter(a.panel.query)
		}
	}
	return false
}

// ── Confirm dialog ────────────────────────────────────────────────────────────

func (a *App) handleConfirmInput(b []byte) bool {
	ch := rune(b[0])
	switch ch {
	case 'y', 'Y':
		if a.confirmAction != nil {
			if err := a.confirmAction(); err != nil {
				a.setStatus("error: " + err.Error())
			} else {
				a.setStatus("done")
			}
		}
		a.mode = modeNormal
		a.confirmMsg = ""
		a.confirmAction = nil
	default:
		a.mode = modeNormal
		a.confirmMsg = ""
		a.confirmAction = nil
		a.setStatus("cancelled")
	}
	return false
}

func (a *App) askConfirm(msg string, action func() error) {
	a.mode = modeConfirm
	a.confirmMsg = msg
	a.confirmAction = action
}

// ── Command bar helpers ───────────────────────────────────────────────────────

func (a *App) enterCommand(prefill string) {
	a.mode = modeCommand
	a.cmdBuf = prefill
	a.statusMsg = ""
}

func (a *App) setStatus(msg string) {
	a.statusMsg = msg
}

// ── Utility ───────────────────────────────────────────────────────────────────

func lastRuneSize(s string) (r rune, size int) {
	if len(s) == 0 {
		return 0, 0
	}
	// scan backward
	for i := len(s); i > 0; {
		r, size = decodeLastRune(s[:i])
		return r, size
	}
	return 0, 0
}

func decodeLastRune(s string) (rune, int) {
	if len(s) == 0 {
		return 0, 0
	}
	// walk backward from end to find rune boundary
	for i := 1; i <= 4 && i <= len(s); i++ {
		r := []rune(s[len(s)-i:])
		if len(r) > 0 {
			return r[0], i
		}
	}
	return rune(s[len(s)-1]), 1
}
