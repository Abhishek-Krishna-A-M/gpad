package index_test

import (
	"testing"
	"time"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/index"
)

func TestNoteEntry_Fields(t *testing.T) {
	entry := &index.NoteEntry{
		RelPath:  "college/math.md",
		Title:    "Mathematics",
		Tags:     []string{"go", "cli"},
		OutLinks: []string{"physics.md", "ideas.md"},
		ModTime:  time.Now(),
	}
	if entry.RelPath != "college/math.md" {
		t.Errorf("RelPath: got %q", entry.RelPath)
	}
	if entry.Title != "Mathematics" {
		t.Errorf("Title: got %q", entry.Title)
	}
	if len(entry.Tags) != 2 {
		t.Errorf("Tags: got %v", entry.Tags)
	}
	if len(entry.OutLinks) != 2 {
		t.Errorf("OutLinks: got %v", entry.OutLinks)
	}
}

func TestCache_Structure(t *testing.T) {
	cache := &index.Cache{
		Notes: map[string]*index.NoteEntry{
			"a.md": {
				RelPath:  "a.md",
				Title:    "Note A",
				Tags:     []string{"tag1"},
				OutLinks: []string{"b"},
			},
			"b.md": {
				RelPath:  "b.md",
				Title:    "Note B",
				Tags:     []string{"tag1", "tag2"},
				OutLinks: []string{},
			},
		},
	}
	if len(cache.Notes) != 2 {
		t.Errorf("expected 2 notes, got %d", len(cache.Notes))
	}
	noteA := cache.Notes["a.md"]
	if noteA == nil {
		t.Fatal("note a.md not found in cache")
	}
	if noteA.Title != "Note A" {
		t.Errorf("title: got %q", noteA.Title)
	}
	if len(noteA.OutLinks) != 1 || noteA.OutLinks[0] != "b" {
		t.Errorf("outlinks: got %v", noteA.OutLinks)
	}
}

func TestLinkEntry_Fields(t *testing.T) {
	entry := &index.LinkEntry{
		OutLinks: []string{"b.md", "c.md"},
		InLinks:  []string{"a.md"},
	}
	if len(entry.OutLinks) != 2 {
		t.Errorf("OutLinks: got %v", entry.OutLinks)
	}
	if len(entry.InLinks) != 1 || entry.InLinks[0] != "a.md" {
		t.Errorf("InLinks: got %v", entry.InLinks)
	}
}

func TestLoad_ReturnsValidEmptyCache(t *testing.T) {
	// When no index.json exists (fresh install / CI), Load must
	// return a non-nil cache with an initialised Notes map.
	cache := index.Load()
	if cache == nil {
		t.Fatal("Load returned nil")
	}
	if cache.Notes == nil {
		t.Error("Load returned cache with nil Notes map")
	}
}
