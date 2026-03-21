// Package links handles [[wikilink]] parsing and the full link graph.
//
// Syntax supported:
//
//	[[note name]]          → resolves to note-name.md anywhere in vault
//	[[folder/note]]        → explicit relative path
//	[[note|display alias]] → alias shown in terminal, link to "note"
//	[[note#heading]]       → heading anchor (stored in index, not jumped)
package links

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)

var wikilinkRe = regexp.MustCompile(`\[\[([^\[\]]+)\]\]`)

// ParseRaw returns the raw inner text of every [[wikilink]] in content.
func ParseRaw(content string) []string {
	matches := wikilinkRe.FindAllStringSubmatch(content, -1)
	var out []string
	for _, m := range matches {
		out = append(out, m[1])
	}
	return out
}

// ParseTargets returns the resolved target name (no alias, no anchor).
func ParseTargets(content string) []string {
	var targets []string
	for _, raw := range ParseRaw(content) {
		t := raw
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
// Returns ("", false) when not found.
func Resolve(target string) (string, bool) {
	notesRoot := storage.NotesDir()
	if !strings.HasSuffix(target, ".md") {
		target += ".md"
	}

	// 1. Exact relative path
	explicit := filepath.Join(notesRoot, target)
	if _, err := os.Stat(explicit); err == nil {
		return explicit, true
	}

	// 2. Walk vault for case-insensitive base name match
	wantBase := strings.ToLower(filepath.Base(target))
	var found string
	_ = filepath.Walk(notesRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.Contains(path, "/.git/") {
			return filepath.SkipDir
		}
		if strings.ToLower(info.Name()) == wantBase {
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

// IndexEntry holds outbound and inbound links for a note.
type IndexEntry struct {
	OutLinks []string // relative paths this note links TO
	InLinks  []string // relative paths that link TO this note
}

// BuildIndex walks the vault and returns the full bidirectional link graph.
// Keys are relative paths ("college/math.md").
func BuildIndex() map[string]*IndexEntry {
	notesRoot := storage.NotesDir()
	index := map[string]*IndexEntry{}

	_ = filepath.Walk(notesRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.Contains(path, "/.git/") {
			return filepath.SkipDir
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		rel, _ := filepath.Rel(notesRoot, path)
		if index[rel] == nil {
			index[rel] = &IndexEntry{}
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		for _, target := range ParseTargets(string(data)) {
			abs, ok := Resolve(target)
			if !ok {
				continue
			}
			targetRel, _ := filepath.Rel(notesRoot, abs)

			// deduplicate outlinks
			if !contains(index[rel].OutLinks, targetRel) {
				index[rel].OutLinks = append(index[rel].OutLinks, targetRel)
			}

			if index[targetRel] == nil {
				index[targetRel] = &IndexEntry{}
			}
			if !contains(index[targetRel].InLinks, rel) {
				index[targetRel].InLinks = append(index[targetRel].InLinks, rel)
			}
		}
		return nil
	})

	return index
}

// Backlinks returns all notes linking TO targetRel.
func Backlinks(targetRel string) []string {
	idx := BuildIndex()
	if e, ok := idx[targetRel]; ok {
		return e.InLinks
	}
	return nil
}

// Outlinks returns all notes linked FROM sourceRel.
func Outlinks(sourceRel string) []string {
	idx := BuildIndex()
	if e, ok := idx[sourceRel]; ok {
		return e.OutLinks
	}
	return nil
}

// ReplaceForDisplay converts [[target]] → "→ target" for terminal rendering.
func ReplaceForDisplay(content string) string {
	return wikilinkRe.ReplaceAllStringFunc(content, func(match string) string {
		inner := match[2 : len(match)-2]
		// use alias if present
		if idx := strings.Index(inner, "|"); idx != -1 {
			return Cyan + "→ " + strings.TrimSpace(inner[idx+1:]) + Reset
		}
		// strip heading anchor
		if idx := strings.Index(inner, "#"); idx != -1 {
			inner = inner[:idx]
		}
		return Cyan + "→ " + strings.TrimSpace(inner) + Reset
	})
}

// ANSI helpers used by ReplaceForDisplay.
const (
	Cyan  = "\033[96m"
	Reset = "\033[0m"
)

func contains(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}
