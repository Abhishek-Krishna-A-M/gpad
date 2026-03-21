// Package daily manages daily notes: one per calendar day in notes/daily/.
package daily

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/editor"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/templates"
)

// Open opens (or creates) today's daily note.
func Open() error {
	return OpenDate(time.Now())
}

// OpenDate opens the daily note for a specific date.
func OpenDate(t time.Time) error {
	path := pathFor(t)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := create(t); err != nil {
			return err
		}
	}

	return editor.Open(path)
}

// pathFor returns the absolute path for a given day's note.
func pathFor(t time.Time) string {
	filename := t.Format("2006-01-02") + ".md"
	return filepath.Join(storage.DailyDir(), filename)
}

// RelPath returns the vault-relative path for today's note.
func RelPath() string {
	t := time.Now()
	filename := t.Format("2006-01-02") + ".md"
	return filepath.Join("daily", filename)
}

func create(t time.Time) error {
	path := pathFor(t)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	// try daily template first
	content, err := templates.Apply("daily", t.Format("2006-01-02"))
	if err != nil {
		// fallback: bare structure
		content = fmt.Sprintf(`---
title: %s
date: %s
tags: [daily]
---

# %s

## Tasks

- [ ] 

## Notes

## Done

`, t.Format("2006-01-02"), t.Format("2006-01-02"), t.Format("Monday, January 2 2006"))
	}

	return os.WriteFile(path, []byte(content), 0644)
}

// Exists reports whether the daily note for t already exists.
func Exists(t time.Time) bool {
	_, err := os.Stat(pathFor(t))
	return err == nil
}

// Yesterday opens yesterday's daily note.
func Yesterday() error {
	return OpenDate(time.Now().AddDate(0, 0, -1))
}

// List returns the last n daily note relative paths, newest first.
func List(n int) []string {
	dir := storage.DailyDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var paths []string
	for i := len(entries) - 1; i >= 0 && len(paths) < n; i-- {
		e := entries[i]
		if !e.IsDir() {
			rel := filepath.Join("daily", e.Name())
			paths = append(paths, rel)
		}
	}
	return paths
}
