package frontmatter_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/frontmatter"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.md")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestParse_WithFrontmatter(t *testing.T) {
	path := writeTemp(t, `---
title: My Note
date: 2026-03-22
tags: [go, cli]
---

# Body here
`)
	meta, body, err := frontmatter.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if meta.Title != "My Note" {
		t.Errorf("title: got %q want %q", meta.Title, "My Note")
	}
	if meta.Date.Format("2006-01-02") != "2026-03-22" {
		t.Errorf("date: got %v", meta.Date)
	}
	if len(meta.Tags) != 2 || meta.Tags[0] != "go" || meta.Tags[1] != "cli" {
		t.Errorf("tags: got %v", meta.Tags)
	}
	if body != "# Body here\n" {
		t.Errorf("body: got %q", body)
	}
}

func TestParse_NoFrontmatter(t *testing.T) {
	path := writeTemp(t, "# Just a note\n\nsome content")
	meta, body, err := frontmatter.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if meta.Title != "" {
		t.Errorf("expected empty title, got %q", meta.Title)
	}
	if body != "# Just a note\n\nsome content" {
		t.Errorf("body mismatch: got %q", body)
	}
}

func TestWrite_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "note.md")

	meta := frontmatter.Meta{
		Title: "Round Trip",
		Date:  time.Date(2026, 3, 22, 0, 0, 0, 0, time.UTC),
		Tags:  []string{"test", "roundtrip"},
	}
	body := "# Round Trip\n\nContent here.\n"

	if err := frontmatter.Write(path, meta, body); err != nil {
		t.Fatal(err)
	}

	meta2, body2, err := frontmatter.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if meta2.Title != meta.Title {
		t.Errorf("title: got %q want %q", meta2.Title, meta.Title)
	}
	if len(meta2.Tags) != 2 {
		t.Errorf("tags: got %v", meta2.Tags)
	}
	if body2 != body {
		t.Errorf("body: got %q want %q", body2, body)
	}
}

func TestAddTag(t *testing.T) {
	path := writeTemp(t, `---
title: Note
date: 2026-03-22
tags: [go]
---

body
`)
	if err := frontmatter.AddTag(path, "cli"); err != nil {
		t.Fatal(err)
	}
	meta, _, _ := frontmatter.Parse(path)
	found := false
	for _, tag := range meta.Tags {
		if tag == "cli" {
			found = true
		}
	}
	if !found {
		t.Errorf("tag 'cli' not found in %v", meta.Tags)
	}
}

func TestAddTag_NoDuplicate(t *testing.T) {
	path := writeTemp(t, `---
title: Note
tags: [go, cli]
---
body
`)
	if err := frontmatter.AddTag(path, "go"); err != nil {
		t.Fatal(err)
	}
	meta, _, _ := frontmatter.Parse(path)
	count := 0
	for _, tag := range meta.Tags {
		if tag == "go" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 occurrence of 'go', got %d in %v", count, meta.Tags)
	}
}

func TestInlineTags(t *testing.T) {
	body := "This is about #programming and #cli tools. Also #go-lang."
	tags := frontmatter.InlineTags(body)
	want := map[string]bool{"programming": true, "cli": true, "go-lang": true}
	for _, tag := range tags {
		if !want[tag] {
			t.Errorf("unexpected tag %q", tag)
		}
		delete(want, tag)
	}
	for remaining := range want {
		t.Errorf("missing tag %q", remaining)
	}
}

func TestRemoveTag(t *testing.T) {
	path := writeTemp(t, `---
title: Note
tags: [go, cli, test]
---
body
`)
	if err := frontmatter.RemoveTag(path, "cli"); err != nil {
		t.Fatal(err)
	}
	meta, _, _ := frontmatter.Parse(path)
	for _, tag := range meta.Tags {
		if tag == "cli" {
			t.Error("tag 'cli' should have been removed")
		}
	}
	if len(meta.Tags) != 2 {
		t.Errorf("expected 2 tags, got %v", meta.Tags)
	}
}
