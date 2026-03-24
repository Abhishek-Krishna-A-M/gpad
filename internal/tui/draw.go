package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/config"
)

// prevW/prevH detect terminal resize between frames.
var prevW, prevH int

// draw builds the entire frame in a strings.Builder and flushes it
// in a single os.Stdout.WriteString call — one syscall per frame.
func (a *App) draw() {
	var sb strings.Builder

	// Hide cursor during draw — prevents visible flicker on st.
	sb.WriteString(aHideCur)

	// On resize: full clear first so old frame doesn't ghost.
	if a.width != prevW || a.height != prevH {
		sb.WriteString(aClearScr)
		prevW, prevH = a.width, a.height
	}

	// Single home — all subsequent writes use mc(row,col).
	sb.WriteString(aHome)

	if a.mode == modeHelp {
		a.drawHelpOverlay(&sb)
	} else if a.panel != nil {
		a.drawPanel(&sb)
	} else {
		a.drawHeader(&sb)
		a.drawMain(&sb)
		a.drawCommandBar(&sb)
		a.drawStatusBar(&sb)
	}

	// Position cursor and restore visibility as the very last operation.
	// This guarantees aShowCur is always the final escape in the buffer.
	switch a.mode {
	case modeCommand:
		// " : " prefix = 3 chars, +1 for the space = col 4 + typed text
		col := 4 + vlen(a.cmdBuf)
		sb.WriteString(mc(a.height-1, col))
	case modeFilter:
		// " / " prefix = 3 chars
		col := 4 + vlen(a.filterStr)
		sb.WriteString(mc(a.height-1, col))
	case modePanelFilter:
		if a.panel != nil {
			col := 4 + vlen(a.panel.query)
			sb.WriteString(mc(a.height-1, col))
		} else {
			sb.WriteString(mc(a.height, 1))
		}
	default:
		// Park cursor at bottom-left — out of the way, not distracting
		sb.WriteString(mc(a.height, 1))
	}
	// Always show cursor as the absolute last byte written
	sb.WriteString(aShowCur)

	os.Stdout.WriteString(sb.String())
}

// ── Header — row 1 ───────────────────────────────────────────────────────────

func (a *App) drawHeader(sb *strings.Builder) {
	sb.WriteString(mc(1, 1))

	noteCount := len(a.flat)
	cwdRel, _ := strings.CutPrefix(a.cwd, a.notesRoot)
	if cwdRel == "" {
		cwdRel = "/"
	}

	title := aBold + aBlue + " gpad" + aReset + aFgDim + cwdRel + aReset
	right := fmt.Sprintf("%s%d notes%s  %s  ", aFgDim, noteCount, aReset, a.gitStatus)

	pad := a.width - vlen(title) - vlen(right)
	if pad < 0 {
		pad = 0
	}
	sb.WriteString(title + rpt(" ", pad) + right + aClearEOL)
}

// ── Main area — rows 2..height-3 ─────────────────────────────────────────────

func (a *App) drawMain(sb *strings.Builder) {
	treeW := a.treeWidth
	prevW := a.width - treeW - 2 // -1 divider -1 padding
	if prevW < 4 {
		prevW = 4
	}
	visH := a.treeHeight()
	preview := a.getPreview()

	for i := 0; i < visH; i++ {
		row := i + 2 // row 1 = header

		// tree column — exactly treeW visible runes
		sb.WriteString(mc(row, 1))
		sb.WriteString(a.renderTreeCell(a.treeScroll+i, treeW))

		// divider — always at treeW+1
		sb.WriteString(mc(row, treeW+1))
		sb.WriteString(aFgMut + "│" + aReset)

		// preview column — starts at treeW+2, truncated to prevW
		sb.WriteString(mc(row, treeW+2))
		line := ""
		if i < len(preview) {
			line = preview[i]
		}
		if vlen(line) > prevW {
			line = truncANSI(line, prevW-1) + "…"
		}
		sb.WriteString(line + aClearEOL)
	}
}

// renderTreeCell returns a string of EXACTLY width visible runes for the tree.
// This guarantees the divider │ always lands on column treeW+1.
func (a *App) renderTreeCell(idx, width int) string {
	if idx < 0 || idx >= len(a.flat) {
		return rpt(" ", width)
	}
	n := a.flat[idx]
	selected := idx == a.cursor

	indent := rpt("  ", n.depth-1)
	icon := "  "
	nameColor := aFg
	suffix := ""

	if n.kind == kindDir {
		nameColor = aBlue
		if n.expanded {
			icon = "▼ "
		} else {
			icon = "▶ "
		}
		suffix = "/"
	} else if config.IsPinned(n.relPath) {
		nameColor = aYellow
		suffix = " ★"
	}

	// measure plain label to compute padding
	plainLabel := indent + icon + n.name + suffix
	plainLen := len([]rune(plainLabel))

	// truncate name if too long
	if plainLen > width-1 {
		maxName := width - 1 - len([]rune(indent)) - len([]rune(icon)) - len([]rune(suffix))
		if maxName < 1 {
			maxName = 1
		}
		nr := []rune(n.name)
		if len(nr) > maxName {
			n = &treeNode{name: string(nr[:maxName-1]) + "…", relPath: n.relPath,
				kind: n.kind, depth: n.depth, expanded: n.expanded}
		}
		plainLabel = indent + icon + n.name + suffix
		plainLen = len([]rune(plainLabel))
	}

	padSpaces := width - 1 - plainLen
	if padSpaces < 0 {
		padSpaces = 0
	}

	if selected {
		return aRev + " " + indent + icon + nameColor + n.name + suffix + aReset +
			rpt(" ", padSpaces)
	}
	return " " + indent + nameColor + icon + n.name + suffix + aReset +
		rpt(" ", padSpaces)
}

// ── Command bar — row height-1 ────────────────────────────────────────────────

func (a *App) drawCommandBar(sb *strings.Builder) {
	sb.WriteString(mc(a.height-1, 1))

	switch a.mode {
	case modeCommand:
		sb.WriteString(aBlue + " : " + aReset + aFg + a.cmdBuf + aReset)
	case modeFilter:
		sb.WriteString(aYellow + " / " + aReset + aFg + a.filterStr + aReset)
	case modePanelFilter:
		if a.panel != nil {
			sb.WriteString(aYellow + " / " + aReset + aFg + a.panel.query + aReset)
		}
	case modeHelp:
		sb.WriteString(aFgDim + " any key to close help" + aReset)
	case modeConfirm:
		sb.WriteString(aRed + aBold + " ! " + aReset + aFg + a.confirmMsg + aReset)
	default:
		if a.statusMsg != "" {
			sb.WriteString(aFgDim + " " + a.statusMsg + aReset)
		} else {
			sb.WriteString(aFgDim + " : cmd  / filter  F search  T tags  W graph  D daily  ? help  q quit" + aReset)
		}
	}
	sb.WriteString(aClearEOL)
}

// ── Status bar — row height ───────────────────────────────────────────────────

func (a *App) drawStatusBar(sb *strings.Builder) {
	sb.WriteString(mc(a.height, 1))

	modeStr := ""
	switch a.mode {
	case modeNormal:
		modeStr = aBold + aBlue + " NORMAL " + aReset
	case modeCommand:
		modeStr = aBold + aGreen + " COMMAND " + aReset
	case modeFilter:
		modeStr = aBold + aYellow + " FILTER " + aReset
	case modePanelFilter:
		modeStr = aBold + aCyan + " SEARCH " + aReset
	case modeHelp:
		modeStr = aBold + aPurple + " HELP " + aReset
	case modeConfirm:
		modeStr = aBold + aRed + " CONFIRM " + aReset
	}

	path := ""
	if n := a.selected(); n != nil {
		path = aFgDim + n.relPath + aReset
	}
	pos := ""
	if len(a.flat) > 0 {
		pos = fmt.Sprintf("%s%d/%d%s", aFgDim, a.cursor+1, len(a.flat), aReset)
	}

	left := modeStr + "  " + path
	right := pos
	padN := a.width - vlen(left) - vlen(right)
	if padN < 0 {
		padN = 0
	}
	sb.WriteString(left + rpt(" ", padN) + right + aClearEOL)
}

// ── Panel overlay — full screen ───────────────────────────────────────────────

func (a *App) drawPanel(sb *strings.Builder) {
	p := a.panel

	// header row
	sb.WriteString(mc(1, 1))
	title := panelTitle(p.kind)
	query := ""
	if p.query != "" {
		query = "  " + aFgDim + p.query + aReset
	}
	sb.WriteString(aBold+aBlue+title+aReset+query+aClearEOL)

	// divider row
	sb.WriteString(mc(2, 1))
	sb.WriteString(aFgDim+rpt("─", a.width)+aReset+aClearEOL)

	// items
	visH := a.height - 4
	results := p.results
	for i := 0; i < visH; i++ {
		row := i + 3
		sb.WriteString(mc(row, 1))
		idx := p.scroll + i
		if idx >= len(results) {
			sb.WriteString(aClearEOL)
			continue
		}
		item := results[idx]
		selected := idx == p.cursor

		label := pad(item.label, a.width/2)
		hint := ""
		if item.sub != "" {
			hint = "  " + aFgDim + item.sub + aReset
		}
		if vlen(label)+vlen(hint) > a.width-2 {
			hint = ""
		}

		if selected {
			sb.WriteString(aRev + " ▶ " + label + hint + aReset + aClearEOL)
		} else {
			sb.WriteString("   " + label + hint + aClearEOL)
		}
	}

	// command bar
	sb.WriteString(mc(a.height-1, 1))
	if a.mode == modePanelFilter {
		sb.WriteString(aYellow + " / " + aReset + aFg + p.query + aReset + aClearEOL)
	} else {
		count := len(results)
		sb.WriteString(aFgDim + fmt.Sprintf(" %d items  /filter  Enter open  Esc close  q quit", count) + aReset + aClearEOL)
	}

	// status bar
	sb.WriteString(mc(a.height, 1))
	pos := ""
	if len(results) > 0 {
		pos = fmt.Sprintf("%s%d/%d%s", aFgDim, p.cursor+1, len(results), aReset)
	}
	modeStr := aBold + aCyan + " " + panelTitle(p.kind) + " " + aReset
	padN := a.width - vlen(modeStr) - vlen(pos)
	if padN < 0 {
		padN = 0
	}
	sb.WriteString(modeStr + rpt(" ", padN) + pos + aClearEOL)
}

func panelTitle(k panelKind) string {
	switch k {
	case panelSearch:
		return " search"
	case panelTags:
		return " tags"
	case panelGraph:
		return " graph"
	case panelDaily:
		return " daily"
	case panelLinks:
		return " links"
	}
	return " panel"
}

// ── Help overlay — full screen ─────────────────────────────────────────────────

func (a *App) drawHelpOverlay(sb *strings.Builder) {
	sections := helpSections()

	// header
	sb.WriteString(mc(1, 1))
	sb.WriteString(aBold + aBlue + " gpad help" + aReset +
		aFgDim + "  any key to close" + aReset + aClearEOL)

	sb.WriteString(mc(2, 1))
	sb.WriteString(aFgDim + strings.Repeat("─", a.width) + aReset + aClearEOL)

	row := 3
	maxRow := a.height - 1
	col := 1
	colWidth := a.width / 2

	for _, sec := range sections {
		if row >= maxRow {
			break
		}
		// section heading
		sb.WriteString(mc(row, col))
		sb.WriteString(aBold + aYellow + " " + sec.title + aReset + aClearEOL)
		row++

		for _, line := range sec.lines {
			if row >= maxRow {
				break
			}
			// two-column layout: if we're past halfway down, start right col
			if row > (maxRow/2)+2 && col == 1 {
				col = colWidth + 1
				row = 3
			}
			sb.WriteString(mc(row, col))
			entry := fmt.Sprintf("  %s%-12s%s %s",
				aCyan, line[0], aReset+aFgDim, line[1]+aReset)
			if vlen(entry) > colWidth-1 {
				entry = truncANSI(entry, colWidth-2)
			}
			sb.WriteString(entry + aClearEOL)
			row++
		}
		row++ // gap between sections
	}

	// statusbar
	sb.WriteString(mc(a.height, 1))
	sb.WriteString(aBold + aPurple + " HELP " + aReset + aClearEOL)
}

type helpSection struct {
	title string
	lines [][2]string // [key, description]
}

func helpSections() []helpSection {
	return []helpSection{
		{
			title: "Navigation",
			lines: [][2]string{
				{"j / ↓", "move down"},
				{"k / ↑", "move up"},
				{"g", "jump to top"},
				{"G", "jump to bottom"},
				{"Enter", "cd into dir / open note"},
				{"l / →", "expand dir in-place"},
				{"h / ←", "collapse / go to parent"},
				{"-", "go up one directory"},
				{"Space", "toggle expand dir"},
			},
		},
		{
			title: "File operations",
			lines: [][2]string{
				{"o", "open in editor"},
				{"v", "view rendered (pager)"},
				{"n", "new note in current dir"},
				{"d / x", "delete (asks confirm)"},
				{"r", "rename (prefills :mv)"},
				{"y", "yank path"},
				{"p", "paste yanked note"},
				{"P", "toggle pin  ★"},
			},
		},
		{
			title: "Panels",
			lines: [][2]string{
				{"F", "search panel (fuzzy+FTS)"},
				{"T", "tag browser"},
				{"W", "graph panel"},
				{"D", "daily notes list"},
				{"L", "links for selected note"},
			},
		},
		{
			title: "Panel keys",
			lines: [][2]string{
				{"j / k", "navigate"},
				{"/", "filter items"},
				{"Enter", "open selected"},
				{"Esc / q", "close panel"},
			},
		},
		{
			title: "General",
			lines: [][2]string{
				{"s", "git sync"},
				{"t", "open today's daily note"},
				{":", "command mode"},
				{"/", "filter tree"},
				{"Ctrl+L", "force redraw"},
				{"q / Esc", "quit"},
			},
		},
		{
			title: "Commands  (press : to enter)",
			lines: [][2]string{
				{":new <n>", "create note in current dir"},
				{":new <n> -t <t>", "create from template"},
				{":mv <src> <dst>", "move / rename"},
				{":cp <src> <dst>", "copy"},
				{":rm <note>", "delete"},
				{":mkdir <dir>", "create directory"},
				{":find <query>", "search panel with results"},
				{":tags", "tag browser panel"},
				{":graph", "graph panel"},
				{":daily", "daily notes panel"},
				{":links", "links panel"},
				{":today", "open today's daily note"},
				{":sync", "pull + push"},
				{":git init <url>", "connect git remote"},
				{":git status", "show remote + autopush"},
				{":config editor <e>", "set editor"},
				{":config autopush on", "enable auto-push"},
				{":template list", "list templates"},
				{":template new <n>", "create template"},
				{":template edit <n>", "edit template"},
				{":keybinds", "edit ~/.gpad/keybinds.json"},
				{":q / :quit", "quit"},
			},
		},
		{
			title: "Wikilinks  (write in notes)",
			lines: [][2]string{
				{"[[note]]", "link to any note by name"},
				{"[[folder/note]]", "explicit path"},
				{"[[note|alias]]", "display alias"},
				{"[[note#heading]]", "heading anchor"},
				{"#tag", "inline tag (also indexed)"},
			},
		},
	}
}
