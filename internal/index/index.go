// Package index maintains a persistent JSON cache of the vault's
// link graph, tag index, and note titles.
//
// The cache is stored in ~/.gpad/index.json and is rebuilt only for
// notes whose mtime has changed since the last run — O(changed) not O(vault).
package index

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/frontmatter"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)

// NoteEntry is the cached record for a single note.
type NoteEntry struct {
	RelPath  string    `json:"rel_path"`
	Title    string    `json:"title"`
	Tags     []string  `json:"tags"`
	OutLinks []string  `json:"out_links"`
	ModTime  time.Time `json:"mod_time"`
}

// Cache is the full persisted index.
type Cache struct {
	Notes   map[string]*NoteEntry `json:"notes"`   // keyed by rel path
	Built   time.Time             `json:"built"`
}

// Load reads the cache from disk. Returns an empty cache on any error.
func Load() *Cache {
	data, err := os.ReadFile(storage.IndexPath())
	if err != nil {
		return &Cache{Notes: map[string]*NoteEntry{}}
	}
	var c Cache
	if err := json.Unmarshal(data, &c); err != nil {
		return &Cache{Notes: map[string]*NoteEntry{}}
	}
	if c.Notes == nil {
		c.Notes = map[string]*NoteEntry{}
	}
	return &c
}

// Save writes the cache to disk.
func (c *Cache) Save() error {
	c.Built = time.Now()
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(storage.IndexPath(), b, 0644)
}

// Refresh walks the vault and updates only notes whose mtime has changed.
// Deleted notes are pruned. Returns true if anything changed.
func Refresh() (*Cache, bool, error) {
	cache := Load()
	notesRoot := storage.NotesDir()
	changed := false

	// track which paths still exist
	seen := map[string]bool{}

	err := filepath.Walk(notesRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		rel, _ := filepath.Rel(notesRoot, path)
		seen[rel] = true

		existing, ok := cache.Notes[rel]
		if ok && !info.ModTime().After(existing.ModTime) {
			// unchanged — skip
			return nil
		}

		// parse the note
		entry, err := parseNote(path, rel, info.ModTime())
		if err != nil {
			return nil
		}
		cache.Notes[rel] = entry
		changed = true
		return nil
	})
	if err != nil {
		return cache, changed, err
	}

	// prune deleted notes
	for rel := range cache.Notes {
		if !seen[rel] {
			delete(cache.Notes, rel)
			changed = true
		}
	}

	if changed {
		_ = cache.Save()
	}
	return cache, changed, nil
}

// parseNote extracts title, tags, and outlinks from a note file.
func parseNote(absPath, relPath string, modTime time.Time) (*NoteEntry, error) {
	meta, body, err := frontmatter.Parse(absPath)
	if err != nil {
		return nil, err
	}

	title := meta.Title
	if title == "" {
		title = strings.TrimSuffix(filepath.Base(relPath), ".md")
	}

	// collect tags from frontmatter + inline
	tagSet := map[string]bool{}
	for _, t := range meta.Tags {
		tagSet[strings.ToLower(t)] = true
	}
	for _, t := range frontmatter.InlineTags(body) {
		tagSet[t] = true
	}
	tags := make([]string, 0, len(tagSet))
	for t := range tagSet {
		tags = append(tags, t)
	}

	// extract wikilink targets
	outLinks := extractLinks(body)

	return &NoteEntry{
		RelPath:  relPath,
		Title:    title,
		Tags:     tags,
		OutLinks: outLinks,
		ModTime:  modTime,
	}, nil
}

// wikilinkRe matches [[...]] — inline to avoid import cycle with links package.
var wikilinkPrefixLen = 2

func extractLinks(content string) []string {
	var targets []string
	seen := map[string]bool{}
	i := 0
	for i < len(content)-3 {
		if content[i] == '[' && content[i+1] == '[' {
			end := strings.Index(content[i+2:], "]]")
			if end >= 0 {
				raw := content[i+2 : i+2+end]
				// strip alias and anchor
				if idx := strings.Index(raw, "|"); idx != -1 {
					raw = raw[:idx]
				}
				if idx := strings.Index(raw, "#"); idx != -1 {
					raw = raw[:idx]
				}
				t := strings.TrimSpace(raw)
				if t != "" && !seen[t] {
					seen[t] = true
					targets = append(targets, t)
				}
				i += 2 + end + 2
				continue
			}
		}
		i++
	}
	return targets
}

// AllNotes returns every NoteEntry in the cache, refreshed if needed.
func AllNotes() []*NoteEntry {
	cache, _, _ := Refresh()
	notes := make([]*NoteEntry, 0, len(cache.Notes))
	for _, e := range cache.Notes {
		notes = append(notes, e)
	}
	return notes
}

// TagIndex returns tag → []relPath from the cache.
func TagIndex() map[string][]string {
	cache, _, _ := Refresh()
	idx := map[string][]string{}
	for rel, e := range cache.Notes {
		for _, t := range e.Tags {
			idx[t] = append(idx[t], rel)
		}
	}
	return idx
}

// LinkIndex returns rel → (outLinks, inLinks) from the cache.
type LinkEntry struct {
	OutLinks []string
	InLinks  []string
}

func LinkIndex() map[string]*LinkEntry {
	cache, _, _ := Refresh()
	notesRoot := storage.NotesDir()
	idx := map[string]*LinkEntry{}

	for rel := range cache.Notes {
		if idx[rel] == nil {
			idx[rel] = &LinkEntry{}
		}
	}

	for rel, e := range cache.Notes {
		for _, target := range e.OutLinks {
			// resolve target name to a rel path
			resolved := resolveTarget(target, notesRoot, cache)
			if resolved == "" {
				continue
			}
			if idx[rel] == nil {
				idx[rel] = &LinkEntry{}
			}
			if !contains(idx[rel].OutLinks, resolved) {
				idx[rel].OutLinks = append(idx[rel].OutLinks, resolved)
			}
			if idx[resolved] == nil {
				idx[resolved] = &LinkEntry{}
			}
			if !contains(idx[resolved].InLinks, rel) {
				idx[resolved].InLinks = append(idx[resolved].InLinks, rel)
			}
		}
	}
	return idx
}

func resolveTarget(target, notesRoot string, cache *Cache) string {
	if !strings.HasSuffix(target, ".md") {
		target += ".md"
	}
	// exact match
	if _, ok := cache.Notes[target]; ok {
		return target
	}
	// case-insensitive base name match
	want := strings.ToLower(filepath.Base(target))
	for rel := range cache.Notes {
		if strings.ToLower(filepath.Base(rel)) == want {
			return rel
		}
	}
	// filesystem fallback
	abs := filepath.Join(notesRoot, target)
	if _, err := os.Stat(abs); err == nil {
		rel, _ := filepath.Rel(notesRoot, abs)
		return rel
	}
	return ""
}

func contains(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}
