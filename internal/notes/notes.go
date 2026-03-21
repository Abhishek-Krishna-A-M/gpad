package notes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/editor"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/frontmatter"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/templates"
)

// Create makes a new note (with optional template) and opens it.
func Create(relPath, templateName string) error {
	base := storage.NotesDir()
	full := filepath.Join(base, relPath)

	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return err
	}
	if _, err := os.Stat(full); err == nil {
		return fmt.Errorf("note already exists: %s", relPath)
	}

	title := strings.TrimSuffix(filepath.Base(relPath), ".md")
	// humanize: replace dashes/underscores with spaces
	titleHuman := strings.NewReplacer("-", " ", "_", " ").Replace(title)

	var content string
	if templateName != "" {
		var err error
		content, err = templates.Apply(templateName, titleHuman)
		if err != nil {
			return err
		}
	} else {
		// default minimal note
		content = fmt.Sprintf(`---
title: %s
date: %s
tags: []
---

# %s

`, titleHuman, time.Now().Format("2006-01-02"), titleHuman)
	}

	if err := os.WriteFile(full, []byte(content), 0644); err != nil {
		return err
	}
	return editor.Open(full)
}

// Open opens an existing note, or creates it if absent.
func Open(relPath string) error {
	base := storage.NotesDir()
	full := filepath.Join(base, relPath)

	if _, err := os.Stat(full); os.IsNotExist(err) {
		return Create(relPath, "")
	}

	// ensure every opened note has frontmatter
	title := strings.TrimSuffix(filepath.Base(relPath), ".md")
	_ = frontmatter.EnsureFrontmatter(full, strings.NewReplacer("-", " ", "_", " ").Replace(title))

	return editor.Open(full)
}

// Stats returns word count, line count, and link count for a note.
func Stats(absPath string) (words, lines, linkCount int, err error) {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return
	}
	content := string(data)
	lines = len(strings.Split(content, "\n"))
	words = len(strings.Fields(content))
	// count [[wikilinks]]
	for i := 0; i < len(content)-1; i++ {
		if content[i] == '[' && content[i+1] == '[' {
			linkCount++
		}
	}
	return
}
