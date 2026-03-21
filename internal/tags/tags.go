// Package tags builds a tag → []notes index across the vault.
// Tags are sourced from YAML frontmatter AND inline #hashtags in body text.
package tags

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/frontmatter"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)

// Index maps tag → sorted list of relative note paths.
type Index map[string][]string

// Build walks the entire vault and constructs the tag index.
func Build() Index {
	notesRoot := storage.NotesDir()
	idx := Index{}

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
		meta, body, err := frontmatter.Parse(path)
		if err != nil {
			return nil
		}

		// collect all tags for this note (deduped)
		seen := map[string]bool{}
		for _, t := range meta.Tags {
			seen[strings.ToLower(t)] = true
		}
		for _, t := range frontmatter.InlineTags(body) {
			seen[t] = true
		}

		for tag := range seen {
			idx[tag] = append(idx[tag], rel)
		}
		return nil
	})

	// sort note lists per tag
	for tag := range idx {
		sort.Strings(idx[tag])
	}
	return idx
}

// NotesForTag returns all relative note paths that carry tag (case-insensitive).
func NotesForTag(tag string) []string {
	tag = strings.ToLower(strings.TrimPrefix(tag, "#"))
	return Build()[tag]
}

// TagsForNote returns all tags on a single note.
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
