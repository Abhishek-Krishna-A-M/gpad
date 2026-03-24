// Package tags builds and queries the tag index across the vault.
// Backed by the persistent index cache — fast on large vaults.
package tags

import (
	"sort"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/frontmatter"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/index"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
	"path/filepath"
)

// Index maps tag → sorted list of relative note paths.
type Index map[string][]string

// Build returns the tag index from the persistent cache.
func Build() Index {
	raw := index.TagIndex()
	idx := make(Index, len(raw))
	for t, notes := range raw {
		sorted := make([]string, len(notes))
		copy(sorted, notes)
		sort.Strings(sorted)
		idx[t] = sorted
	}
	return idx
}

// NotesForTag returns notes tagged with tag (case-insensitive).
func NotesForTag(tag string) []string {
	tag = strings.ToLower(strings.TrimPrefix(tag, "#"))
	return Build()[tag]
}

// TagsForNote returns all tags on a single note (frontmatter + inline).
func TagsForNote(relPath string) []string {
	absPath := filepath.Join(storage.NotesDir(), relPath)
	meta, body, err := frontmatter.Parse(absPath)
	if err != nil {
		return nil
	}
	seen := map[string]bool{}
	var out []string
	add := func(t string) {
		t = strings.ToLower(t)
		if !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
	}
	for _, t := range meta.Tags {
		add(t)
	}
	for _, t := range frontmatter.InlineTags(body) {
		add(t)
	}
	sort.Strings(out)
	return out
}

// AllTags returns every tag in the vault sorted alphabetically.
func AllTags() []string {
	idx := Build()
	all := make([]string, 0, len(idx))
	for t := range idx {
		all = append(all, t)
	}
	sort.Strings(all)
	return all
}
