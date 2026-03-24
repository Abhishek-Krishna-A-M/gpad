package tui

// ── Cursor movement ───────────────────────────────────────────────────────────

func (a *App) moveUp() {
	if a.cursor > 0 {
		a.cursor--
		a.previewCache = ""
		a.scrollIntoView()
	}
}

func (a *App) moveDown() {
	if a.cursor < len(a.flat)-1 {
		a.cursor++
		a.previewCache = ""
		a.scrollIntoView()
	}
}

func (a *App) moveTop() {
	if len(a.flat) == 0 {
		return
	}
	a.cursor = 0
	a.treeScroll = 0
	a.previewCache = ""
}

func (a *App) moveBottom() {
	if len(a.flat) == 0 {
		return
	}
	a.cursor = len(a.flat) - 1
	a.previewCache = ""
	a.scrollIntoView()
}

func (a *App) scrollIntoView() {
	vh := a.treeHeight()
	if vh <= 0 {
		return
	}
	if a.cursor < a.treeScroll {
		a.treeScroll = a.cursor
	}
	if a.cursor >= a.treeScroll+vh {
		a.treeScroll = a.cursor - vh + 1
	}
}

func (a *App) treeHeight() int {
	// header(1) + cmdbar(1) + statusbar(1) + 1 padding = 4 reserved rows
	h := a.height - 4
	if h < 1 {
		h = 1
	}
	return h
}

// selected returns the currently highlighted treeNode, or nil.
func (a *App) selected() *treeNode {
	if len(a.flat) == 0 || a.cursor < 0 || a.cursor >= len(a.flat) {
		return nil
	}
	return a.flat[a.cursor]
}

// ── Enter key — lf-style ──────────────────────────────────────────────────────

// handleEnter implements lf Enter semantics:
//   - on a dir  → cd into it (tree view replaced with dir contents)
//   - on a note → open in editor
func (a *App) handleEnter() {
	n := a.selected()
	if n == nil {
		return
	}
	if n.kind == kindDir {
		a.cdInto(n)
	} else {
		a.openSelected()
	}
}

// handleExpand implements Space semantics:
//   - on a dir  → expand/collapse in-place (tree style)
//   - on a note → open in editor (same as Enter)
func (a *App) handleExpand() {
	n := a.selected()
	if n == nil {
		return
	}
	if n.kind == kindDir {
		a.expandInPlace(n)
	} else {
		a.openSelected()
	}
}
