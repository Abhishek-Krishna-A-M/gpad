package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/daily"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/links"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/search"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/tags"
)

// ── Panel type ────────────────────────────────────────────────────────────────

type panelKind int

const (
	panelNone   panelKind = iota
	panelSearch
	panelTags
	panelGraph
	panelDaily
	panelLinks
)

type PanelState struct {
	kind    panelKind
	items   []panelItem  // full unfiltered list
	cursor  int
	scroll  int
	query   string
	results []panelItem  // filtered/searched results
}

type panelItem struct {
	label string
	sub   string
	value string // rel path or tag name
}

// ── Open / close ──────────────────────────────────────────────────────────────

func (a *App) openPanel(kind panelKind) {
	a.panel = &PanelState{kind: kind}
	switch kind {
	case panelSearch:
		a.loadSearchResults("")
		a.mode = modePanelFilter
	case panelTags:
		a.loadTagPanel()
	case panelGraph:
		a.loadGraphPanel()
	case panelDaily:
		a.loadDailyPanel()
	case panelLinks:
		a.loadLinksPanel()
	}
}

func (a *App) closePanel() {
	a.panel = nil
	a.mode = modeNormal
}

// ── Search panel — uses search.FullText + search.Fuzzy ───────────────────────

func (a *App) loadSearchResults(query string) {
	if a.panel == nil {
		return
	}
	a.panel.query = query
	a.panel.cursor = 0
	a.panel.scroll = 0

	if strings.TrimSpace(query) == "" {
		all := search.AllNotePaths()
		items := make([]panelItem, len(all))
		for i, p := range all {
			items[i] = panelItem{label: p, value: p}
		}
		a.panel.results = items
		return
	}

	seen := map[string]bool{}
	var items []panelItem

	// full-text body search — match count scored
	for _, r := range search.FullText(query) {
		if seen[r.RelPath] { continue }
		seen[r.RelPath] = true
		sub := r.Excerpt
		if sub == "" { sub = fmt.Sprintf("×%d", r.Score) }
		items = append(items, panelItem{label: r.RelPath, sub: sub, value: r.RelPath})
	}
	// fuzzy title match — catches notes not in body results
	for _, r := range search.Fuzzy(query, search.AllNotePaths()) {
		if seen[r.RelPath] { continue }
		seen[r.RelPath] = true
		items = append(items, panelItem{label: r.RelPath, value: r.RelPath})
	}
	a.panel.results = items
}

// ── Tag panel — uses tags.AllTags + tags.Build ────────────────────────────────

func (a *App) loadTagPanel() {
	all := tags.AllTags()
	idx := tags.Build()
	items := make([]panelItem, len(all))
	for i, t := range all {
		items[i] = panelItem{
			label: "#" + t,
			sub:   fmt.Sprintf("%d notes", len(idx[t])),
			value: t,
		}
	}
	a.panel.items = items
	a.panel.results = items
}

// tagDrilldown — Enter on a tag shows its notes via tags.NotesForTag
func (a *App) tagDrilldown(tag string) {
	noteList := tags.NotesForTag(tag)
	items := make([]panelItem, len(noteList))
	for i, n := range noteList {
		items[i] = panelItem{label: n, value: n}
	}
	a.panel.items = items
	a.panel.results = items
	a.panel.cursor = 0
	a.panel.scroll = 0
	a.panel.query = "#" + tag
}

// ── Graph panel — uses links.BuildIndex ──────────────────────────────────────

func (a *App) loadGraphPanel() {
	idx := links.BuildIndex()

	type entry struct {
		rel      string
		outCount int
		inCount  int
	}
	var entries []entry
	for rel, e := range idx {
		if len(e.OutLinks)+len(e.InLinks) > 0 {
			entries = append(entries, entry{rel, len(e.OutLinks), len(e.InLinks)})
		}
	}
	// sort by total connections descending — most connected notes first
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].outCount+entries[i].inCount > entries[j].outCount+entries[j].inCount
	})

	items := make([]panelItem, len(entries))
	for i, e := range entries {
		items[i] = panelItem{
			label: e.rel,
			sub:   fmt.Sprintf("→%d ←%d", e.outCount, e.inCount),
			value: e.rel,
		}
	}
	a.panel.items = items
	a.panel.results = items
}

// ── Daily panel — uses daily.List ────────────────────────────────────────────

func (a *App) loadDailyPanel() {
	// daily.List returns the last 90 daily notes newest-first
	relPaths := daily.List(90)
	if len(relPaths) == 0 {
		a.panel.results = []panelItem{{label: "no daily notes yet — press t to create today's", value: ""}}
		return
	}

	todayRel := daily.RelPath()
	items := make([]panelItem, len(relPaths))
	for i, rel := range relPaths {
		date := strings.TrimSuffix(filepath.Base(rel), ".md")
		sub := ""
		if rel == todayRel {
			sub = aGreen+"today"+aReset
		}
		items[i] = panelItem{label: date, sub: sub, value: rel}
	}
	a.panel.items = items
	a.panel.results = items
}

// ── Links panel — uses links.Outlinks + links.Backlinks ──────────────────────

func (a *App) loadLinksPanel() {
	n := a.selected()
	if n == nil || n.kind != kindNote {
		a.panel.results = []panelItem{{label: "select a note first", value: ""}}
		return
	}

	out := links.Outlinks(n.relPath)
	back := links.Backlinks(n.relPath)

	var items []panelItem
	for _, o := range out {
		items = append(items, panelItem{
			label: o,
			sub:   aCyan+"→ outlink"+aReset,
			value: o,
		})
	}
	for _, b := range back {
		items = append(items, panelItem{
			label: b,
			sub:   aGreen+"← backlink"+aReset,
			value: b,
		})
	}
	if len(items) == 0 {
		items = []panelItem{{label: "no links — add [[note name]] to connect notes", value: ""}}
	}
	a.panel.items = items
	a.panel.results = items
	a.panel.query = n.name
}

// ── Panel navigation ──────────────────────────────────────────────────────────

func (a *App) panelMoveUp() {
	if a.panel.cursor > 0 {
		a.panel.cursor--
		a.panelScrollIntoView()
	}
}

func (a *App) panelMoveDown() {
	if a.panel.cursor < len(a.panel.results)-1 {
		a.panel.cursor++
		a.panelScrollIntoView()
	}
}

func (a *App) panelScrollIntoView() {
	vh := a.treeHeight()
	if a.panel.cursor < a.panel.scroll {
		a.panel.scroll = a.panel.cursor
	}
	if a.panel.cursor >= a.panel.scroll+vh {
		a.panel.scroll = a.panel.cursor - vh + 1
	}
}

func (a *App) panelEnter() {
	if a.panel == nil || len(a.panel.results) == 0 { return }
	if a.panel.cursor >= len(a.panel.results) { return }

	item := a.panel.results[a.panel.cursor]
	if item.value == "" { return }

	switch a.panel.kind {
	case panelSearch, panelLinks, panelGraph, panelDaily:
		absPath := filepath.Join(a.notesRoot, item.value)
		a.closePanel()
		a.openAbsPath(absPath, item.value)

	case panelTags:
		if strings.HasPrefix(item.label, "#") {
			// on a tag — drill into its notes
			a.tagDrilldown(item.value)
		} else {
			// on a note inside a tag — open it
			absPath := filepath.Join(a.notesRoot, item.value)
			a.closePanel()
			a.openAbsPath(absPath, item.value)
		}
	}
}

func (a *App) panelFilter(query string) {
	a.panel.query = query
	a.panel.cursor = 0
	a.panel.scroll = 0

	if a.panel.kind == panelSearch {
		a.loadSearchResults(query)
		return
	}

	if query == "" {
		a.panel.results = a.panel.items
		return
	}
	pat := strings.ToLower(query)
	var out []panelItem
	for _, it := range a.panel.items {
		if fuzzyMatch(pat, strings.ToLower(it.label)) {
			out = append(out, it)
		}
	}
	a.panel.results = out
}

func fuzzyMatch(pattern, target string) bool {
	pi := 0
	pr := []rune(pattern)
	for _, tc := range target {
		if pi < len(pr) && tc == pr[pi] { pi++ }
	}
	return pi == len(pr)
}

// runGraphExternal is a fallback that spawns gpad graph in the pager.
func (a *App) runGraphExternal() {
	a.suspendAndRun(func() error {
		cmd := exec.Command("gpad", "graph")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("tip: press W to open graph panel inside TUI")
		}
		fmt.Print("\npress any key…")
		buf := make([]byte, 1)
		os.Stdin.Read(buf)
		return nil
	})
}
