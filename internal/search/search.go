// Package search provides full-text and fuzzy search backed by the index cache.
package search

import (
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/index"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
	"path/filepath"
)

// Result is a single search hit.
type Result struct {
	RelPath string
	Title   string
	Excerpt string
	Score   int
}

// FullText searches every note's body for query (case-insensitive).
// Uses the index for the note list but reads files for body content.
func FullText(query string) []Result {
	if query == "" {
		return nil
	}
	notesRoot := storage.NotesDir()
	q := strings.ToLower(query)
	var results []Result

	notes := index.AllNotes()
	for _, n := range notes {
		absPath := filepath.Join(notesRoot, n.RelPath)
		data, err := os.ReadFile(absPath)
		if err != nil {
			continue
		}
		content := string(data)
		lower := strings.ToLower(content)
		count := strings.Count(lower, q)
		if count == 0 {
			continue
		}
		results = append(results, Result{
			RelPath: n.RelPath,
			Title:   n.Title,
			Excerpt: extractExcerpt(content, query),
			Score:   count,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results
}

// Fuzzy filters notes by a fuzzy pattern against their title and path.
func Fuzzy(pattern string, notePaths []string) []Result {
	if pattern == "" {
		out := make([]Result, 0, len(notePaths))
		for _, n := range notePaths {
			out = append(out, Result{RelPath: n})
		}
		return out
	}
	p := strings.ToLower(pattern)
	var results []Result
	for _, note := range notePaths {
		score := fuzzyScore(p, strings.ToLower(note))
		if score > 0 {
			results = append(results, Result{RelPath: note, Score: score})
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results
}

func fuzzyScore(pattern, target string) int {
	if pattern == "" {
		return 1
	}
	pi := 0
	score := 0
	consecutive := 0
	pr := []rune(pattern)
	for ti, tc := range target {
		if pi < len(pr) && unicode.ToLower(tc) == unicode.ToLower(pr[pi]) {
			pi++
			consecutive++
			score += consecutive * 2
			if ti == 0 {
				score += 5
			}
		} else {
			consecutive = 0
		}
	}
	if pi != len(pr) {
		return 0
	}
	if target == pattern {
		score += 100
	}
	return score
}

func extractExcerpt(content, query string) string {
	lines := strings.Split(content, "\n")
	q := strings.ToLower(query)
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), q) {
			trimmed := strings.TrimSpace(line)
			// skip frontmatter lines
			if strings.HasPrefix(trimmed, "---") {
				continue
			}
			if strings.Contains(trimmed, ":") && len(trimmed) < 40 {
				continue
			}
			if len(trimmed) > 120 {
				idx := strings.Index(strings.ToLower(trimmed), q)
				start := idx - 30
				if start < 0 {
					start = 0
				}
				end := idx + len(query) + 60
				if end > len(trimmed) {
					end = len(trimmed)
				}
				trimmed = "…" + trimmed[start:end] + "…"
			}
			return trimmed
		}
	}
	return ""
}

// AllNotePaths returns relative paths of all notes from the index cache.
func AllNotePaths() []string {
	notes := index.AllNotes()
	paths := make([]string, 0, len(notes))
	for _, n := range notes {
		paths = append(paths, n.RelPath)
	}
	sort.Strings(paths)
	return paths
}
