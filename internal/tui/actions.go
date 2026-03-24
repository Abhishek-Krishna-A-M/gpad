package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/config"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/core"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/daily"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/editor"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/notes"
)

// ── Editor open ───────────────────────────────────────────────────────────────

func (a *App) openSelected() {
	n := a.selected()
	if n == nil || n.kind != kindNote {
		a.setStatus("select a note to open")
		return
	}
	a.openAbsPath(n.absPath, n.relPath)
}

// openAbsPath suspends the TUI, opens absPath in the user's editor,
// then resumes instantly. Git push runs in a background goroutine so
// there is zero lag when returning to the TUI.
func (a *App) openAbsPath(absPath, relPath string) {
	a.suspendAndRun(func() error {
		return editor.Open(absPath)
	})
	a.previewCache = ""
	_ = a.buildTree()
	// async push — TUI is already redrawn before git does any network I/O
	go func() {
		core.AutoSave("update " + relPath)
	}()
	a.setStatus("↑ pushing…")
	go func() {
		time.Sleep(3 * time.Second)
		a.setStatus("")
	}()
}

// ── View rendered ─────────────────────────────────────────────────────────────

func (a *App) viewSelected() {
	n := a.selected()
	if n == nil || n.kind != kindNote {
		a.setStatus("select a note to view")
		return
	}
	// viewer.View renders markdown + backlinks + stats then pages with less.
	// We suspend the TUI around it — same pattern as openAbsPath.
	a.suspendAndRun(func() error {
		// pipe rendered output through less ourselves so we don't get the
		// "edit?" prompt that viewer.View normally shows
		lines := renderPreview(n.absPath, 9999)
		cmd := exec.Command("less", "-R", "--quit-if-one-screen")
		cmd.Stdin = strings.NewReader(strings.Join(lines, "\n"))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	})
}

// ── Delete ────────────────────────────────────────────────────────────────────

func (a *App) deleteSelected() {
	n := a.selected()
	if n == nil {
		return
	}
	a.askConfirm(
		fmt.Sprintf("delete %s? [y/N]", n.relPath),
		func() error {
			recursive := n.kind == kindDir
			if err := core.Delete([]string{n.relPath}, recursive); err != nil {
				return err
			}
			_ = a.buildTree()
			return nil
		},
	)
}

// ── Rename (prefills command bar with mv) ─────────────────────────────────────

func (a *App) renameSelected() {
	n := a.selected()
	if n == nil {
		return
	}
	a.enterCommand("mv " + n.relPath + " ")
}

// ── Move / Copy ───────────────────────────────────────────────────────────────

func (a *App) runMv(src, dest string) {
	if err := core.Move([]string{src}, dest); err != nil {
		a.setStatus("mv: " + err.Error())
		return
	}
	a.setStatus(fmt.Sprintf("moved %s → %s", src, dest))
	_ = a.buildTree()
}

func (a *App) runCp(src, dest string) {
	if err := core.Copy(src, dest); err != nil {
		a.setStatus("cp: " + err.Error())
		return
	}
	a.setStatus(fmt.Sprintf("copied %s → %s", src, dest))
	_ = a.buildTree()
}

// ── Yank / Paste ──────────────────────────────────────────────────────────────

func (a *App) yankSelected() {
	n := a.selected()
	if n == nil {
		return
	}
	a.yankBuf = n.relPath
	a.setStatus("yanked: " + n.relPath)
}

func (a *App) pasteYanked() {
	if a.yankBuf == "" {
		a.setStatus("nothing yanked — use y first")
		return
	}
	n := a.selected()
	dest := ""
	if n != nil {
		if n.kind == kindDir {
			dest = n.relPath
		} else {
			dest = filepath.Dir(n.relPath)
		}
	}
	if dest == "." || dest == "" {
		dest = filepath.Base(a.yankBuf)
	} else {
		dest = filepath.Join(dest, filepath.Base(a.yankBuf))
	}
	a.runCp(a.yankBuf, dest)
}

// ── New note ──────────────────────────────────────────────────────────────────

func (a *App) newNote() {
	// cwdRel: the current working directory relative to notesRoot.
	// If the user has cd'd into daily/ with Enter, new files land there.
	cwdRel := a.cwdRelative()

	// refine further: if cursor is on a node, use its directory
	n := a.selected()
	dir := cwdRel
	if n != nil {
		if n.kind == kindDir {
			dir = n.relPath
		} else if filepath.Dir(n.relPath) != "." {
			dir = filepath.Dir(n.relPath)
		}
	}
	if dir == "." || dir == "" {
		a.enterCommand("new ")
	} else {
		a.enterCommand("new " + dir + "/")
	}
}

func (a *App) newDir() {
	cwdRel := a.cwdRelative()
	if cwdRel == "" || cwdRel == "." {
		a.enterCommand("mkdir ")
	} else {
		a.enterCommand("mkdir " + cwdRel + "/")
	}
}

// cwdRelative returns a.cwd relative to a.notesRoot.
// Returns "" when we are at the vault root.
func (a *App) cwdRelative() string {
	rel, err := filepath.Rel(a.notesRoot, a.cwd)
	if err != nil || rel == "." {
		return ""
	}
	return rel
}

// createNote uses notes.Create() so templates and frontmatter are handled
// by our existing notes package — no reimplementation.
// name is relative to notesRoot (notes.Create expects that).
func (a *App) createNote(name, tmpl string) {
	if name == "" {
		a.setStatus("usage: new <name.md> [-t template]")
		return
	}
	if !strings.HasSuffix(name, ".md") {
		name += ".md"
	}
	a.suspendAndRun(func() error {
		return notes.Create(name, tmpl)
	})
	a.previewCache = ""
	_ = a.buildTree()
	go func() {
		core.AutoSave("new " + name)
	}()
	a.setStatus("↑ pushing…")
	go func() {
		time.Sleep(3 * time.Second)
		a.setStatus("")
	}()
}

func (a *App) createDir(name string) {
	// resolve relative to cwd, not always notesRoot
	var absPath string
	if filepath.IsAbs(name) {
		absPath = name
	} else {
		// if name already contains a path separator it's explicit
		if strings.ContainsRune(name, filepath.Separator) {
			absPath = filepath.Join(a.notesRoot, name)
		} else {
			absPath = filepath.Join(a.cwd, name)
		}
	}
	if err := os.MkdirAll(absPath, 0755); err != nil {
		a.setStatus("mkdir: " + err.Error())
		return
	}
	rel, _ := filepath.Rel(a.notesRoot, absPath)
	a.setStatus("created " + rel)
	_ = a.buildTree()
}

// ── Pin / Unpin ───────────────────────────────────────────────────────────────

func (a *App) togglePin() {
	n := a.selected()
	if n == nil || n.kind != kindNote {
		a.setStatus("select a note to pin")
		return
	}
	if config.IsPinned(n.relPath) {
		_ = config.Unpin(n.relPath)
		a.setStatus("unpinned " + n.relPath)
	} else {
		_ = config.Pin(n.relPath)
		a.setStatus("★ pinned " + n.relPath)
	}
}

// ── Sync ──────────────────────────────────────────────────────────────────────

func (a *App) runSync() {
	a.setStatus("syncing…")
	a.draw()
	if err := core.Sync(); err != nil {
		a.setStatus("sync error: " + err.Error())
	} else {
		a.gitStatus = a.checkGit()
		a.setStatus("synced")
		_ = a.buildTree()
	}
}

// ── Today's daily note ────────────────────────────────────────────────────────

// openToday uses daily.Open() — consistent with the gpad today command.
// Creates the note from the daily template if it doesn't exist yet.
func (a *App) openToday() {
	a.suspendAndRun(func() error {
		return daily.Open()
	})
	a.previewCache = ""
	_ = a.buildTree()
	go func() {
		core.AutoSave("daily " + daily.RelPath())
	}()
	a.setStatus("↑ pushing…")
	go func() {
		time.Sleep(3 * time.Second)
		a.setStatus("")
	}()
}

// ── Git status ────────────────────────────────────────────────────────────────

func (a *App) checkGit() string {
	out, err := exec.Command("git", "-C", a.notesRoot, "status", "--porcelain").Output()
	if err != nil {
		cfg, _ := config.Load()
		if !cfg.GitEnabled {
			return aFgDim + "no git" + aReset
		}
		return aFgDim + "git err" + aReset
	}
	if len(strings.TrimSpace(string(out))) == 0 {
		return aGreen + "✓" + aReset
	}
	return aYellow + "●" + aReset
}

// ── Suspend terminal for external programs ────────────────────────────────────

func (a *App) suspendAndRun(fn func() error) {
	restoreTerminal(a.fd, a.oldState)
	fmt.Print(aShowCur + aClearScr + aHome)

	if err := fn(); err != nil {
		fmt.Printf("\nerror: %v\npress enter to continue…", err)
		buf := make([]byte, 1)
		os.Stdin.Read(buf)
	}

	if state, err := makeRaw(a.fd); err == nil {
		a.oldState = state
	}
	fmt.Print(aHideCur + aClearScr)
}
