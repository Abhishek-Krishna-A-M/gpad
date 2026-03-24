package tui

import (
	"os"
	"path/filepath"
	"strings"
)

// ── Node types ────────────────────────────────────────────────────────────────

type nodeKind int

const (
	kindDir  nodeKind = iota
	kindNote          // .md file
)

type treeNode struct {
	name     string
	relPath  string // relative to notesRoot
	absPath  string
	kind     nodeKind
	depth    int
	expanded bool
	children []*treeNode
	parent   *treeNode
}

// ── Tree construction ─────────────────────────────────────────────────────────

// buildTree rebuilds the flat visible list from the notes root.
// Preserves cursor if possible — finds the same relPath after rebuild.
func (a *App) buildTree() error {
	// remember what was selected so we can restore after rebuild
	var prevRel string
	if n := a.selected(); n != nil {
		prevRel = n.relPath
	}

	root, err := buildNode(a.cwd, "", 0, a.notesRoot)
	if err != nil {
		return err
	}

	a.treeRoot = root
	a.flat = flatten(root)

	// restore cursor to same node
	if prevRel != "" {
		for i, n := range a.flat {
			if n.relPath == prevRel {
				a.cursor = i
				a.scrollIntoView()
				return nil
			}
		}
	}

	// fallback: clamp cursor
	if a.cursor >= len(a.flat) {
		a.cursor = len(a.flat) - 1
	}
	if a.cursor < 0 {
		a.cursor = 0
	}
	return nil
}

func buildNode(absPath, relPath string, depth int, notesRoot string) (*treeNode, error) {
	entries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, err
	}

	node := &treeNode{
		name:     filepath.Base(absPath),
		relPath:  relPath,
		absPath:  absPath,
		kind:     kindDir,
		depth:    depth,
		expanded: true, // cwd root is always expanded
	}

	var dirs, files []os.DirEntry
	for _, e := range entries {
		name := e.Name()
		if name == ".git" || strings.HasPrefix(name, ".") {
			continue
		}
		if e.IsDir() {
			dirs = append(dirs, e)
		} else if strings.HasSuffix(name, ".md") {
			files = append(files, e)
		}
	}

	for _, e := range dirs {
		childAbs := filepath.Join(absPath, e.Name())
		childRel := e.Name()
		if relPath != "" {
			childRel = filepath.Join(relPath, e.Name())
		}
		child := &treeNode{
			name:    e.Name(),
			relPath: childRel,
			absPath: childAbs,
			kind:    kindDir,
			depth:   depth + 1,
			parent:  node,
		}
		node.children = append(node.children, child)
	}

	for _, e := range files {
		childAbs := filepath.Join(absPath, e.Name())
		childRel := e.Name()
		if relPath != "" {
			childRel = filepath.Join(relPath, e.Name())
		}
		child := &treeNode{
			name:    e.Name(),
			relPath: childRel,
			absPath: childAbs,
			kind:    kindNote,
			depth:   depth + 1,
			parent:  node,
		}
		node.children = append(node.children, child)
	}

	return node, nil
}

// flatten returns the visible nodes in display order.
// Dirs always come before files at each level.
// Only expanded dirs show their children.
func flatten(root *treeNode) []*treeNode {
	var out []*treeNode
	var walk func(n *treeNode)
	walk = func(n *treeNode) {
		for _, child := range n.children {
			out = append(out, child)
			if child.kind == kindDir && child.expanded {
				walk(child)
			}
		}
	}
	walk(root)
	return out
}

// ── lf-style directory navigation ─────────────────────────────────────────────

// cdInto changes the current working directory to the selected dir.
// This is the lf-style Enter on a directory: replaces the tree view
// with the contents of the dir, depth resets to 0.
func (a *App) cdInto(n *treeNode) {
	if n == nil || n.kind != kindDir {
		return
	}
	a.cwdStack = append(a.cwdStack, a.cwd)
	a.cwd = n.absPath
	a.cursor = 0
	a.treeScroll = 0
	a.previewCache = ""
	_ = a.buildTree()
}

// cdUp goes up one directory level (lf: pressing h or - at root of cwd).
func (a *App) cdUp() {
	if len(a.cwdStack) > 0 {
		// pop the stack
		prev := a.cwdStack[len(a.cwdStack)-1]
		a.cwdStack = a.cwdStack[:len(a.cwdStack)-1]
		a.cwd = prev
	} else {
		// not in stack — go to parent dir if not already at notesRoot
		parent := filepath.Dir(a.cwd)
		if parent == a.cwd || parent == filepath.Dir(a.notesRoot) {
			return // already at top
		}
		a.cwd = parent
	}
	a.cursor = 0
	a.treeScroll = 0
	a.previewCache = ""
	_ = a.buildTree()
}

// expandInPlace toggles a dir's expanded state without cd-ing into it.
// This is the Space binding — tree-style expand like nnn.
func (a *App) expandInPlace(n *treeNode) {
	if n == nil || n.kind != kindDir {
		return
	}
	n.expanded = !n.expanded
	prevRel := n.relPath
	a.flat = flatten(a.treeRoot)
	// restore cursor to the toggled dir
	for i, node := range a.flat {
		if node.relPath == prevRel {
			a.cursor = i
			a.scrollIntoView()
			return
		}
	}
}

// collapseOrUp collapses an expanded dir, or cds up if already collapsed/on file.
// h / left arrow behaviour — mirrors lf exactly.
func (a *App) collapseOrUp() {
	n := a.selected()
	if n == nil {
		a.cdUp()
		return
	}
	if n.kind == kindDir && n.expanded {
		// collapse in-place, stay on dir
		n.expanded = false
		prevRel := n.relPath
		a.flat = flatten(a.treeRoot)
		for i, node := range a.flat {
			if node.relPath == prevRel {
				a.cursor = i
				a.scrollIntoView()
				return
			}
		}
		return
	}
	// on a file or collapsed dir → cd up
	a.cdUp()
}
