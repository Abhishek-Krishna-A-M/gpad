// Package links handles [[wikilink]] parsing and the bidirectional link graph.
// The graph is built from the persistent index cache — O(changed notes) not O(vault).
//
// Syntax:
//   [[note name]]          resolves anywhere in vault (case-insensitive)
//   [[folder/note]]        explicit relative path
//   [[note|alias]]         alias displayed, "note" linked
//   [[note#heading]]       heading anchor stored in index
package links

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/index"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)

var wikilinkRe = regexp.MustCompile(`\[\[([^\[\]]+)\]\]`)

// ParseTargets returns the resolved target names from [[wikilinks]] in content.
func ParseTargets(content string) []string {
	matches := wikilinkRe.FindAllStringSubmatch(content, -1)
	var targets []string
	for _, m := range matches {
		t := m[1]
		if idx := strings.Index(t, "|"); idx != -1 {
			t = t[:idx]
		}
		if idx := strings.Index(t, "#"); idx != -1 {
			t = t[:idx]
		}
		targets = append(targets, strings.TrimSpace(t))
	}
	return targets
}

// Resolve turns a wikilink target string into an absolute path.
func Resolve(target string) (string, bool) {
	notesRoot := storage.NotesDir()
	if !strings.HasSuffix(target, ".md") {
		target += ".md"
	}
	explicit := filepath.Join(notesRoot, target)
	if _, err := os.Stat(explicit); err == nil {
		return explicit, true
	}
	want := strings.ToLower(filepath.Base(target))
	var found string
	_ = filepath.Walk(notesRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.Contains(path, "/.git/") {
			return filepath.SkipDir
		}
		if strings.ToLower(info.Name()) == want {
			found = path
			return filepath.SkipAll
		}
		return nil
	})
	if found != "" {
		return found, true
	}
	return "", false
}

// Backlinks returns all notes linking TO targetRel, using the index cache.
func Backlinks(targetRel string) []string {
	idx := index.LinkIndex()
	if e, ok := idx[targetRel]; ok {
		return e.InLinks
	}
	return nil
}

// Outlinks returns all notes linked FROM sourceRel, using the index cache.
func Outlinks(sourceRel string) []string {
	idx := index.LinkIndex()
	if e, ok := idx[sourceRel]; ok {
		return e.OutLinks
	}
	return nil
}

// BuildIndex returns the full bidirectional link graph from the index cache.
func BuildIndex() map[string]*IndexEntry {
	cached := index.LinkIndex()
	out := make(map[string]*IndexEntry, len(cached))
	for k, v := range cached {
		out[k] = &IndexEntry{OutLinks: v.OutLinks, InLinks: v.InLinks}
	}
	return out
}

// IndexEntry holds outbound and inbound links for a note.
type IndexEntry struct {
	OutLinks []string
	InLinks  []string
}

// ReplaceForDisplay converts [[target]] to a coloured terminal label.
func ReplaceForDisplay(content string) string {
	return wikilinkRe.ReplaceAllStringFunc(content, func(match string) string {
		inner := match[2 : len(match)-2]
		if idx := strings.Index(inner, "|"); idx != -1 {
			return Cyan + "→ " + strings.TrimSpace(inner[idx+1:]) + Reset
		}
		if idx := strings.Index(inner, "#"); idx != -1 {
			inner = inner[:idx]
		}
		return Cyan + "→ " + strings.TrimSpace(inner) + Reset
	})
}

const (
	Cyan  = "\033[96m"
	Reset = "\033[0m"
)
