// Package search provides full-text search and fuzzy title filtering
// across the gpad vault without any external dependencies.
package search

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/frontmatter"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)

// Result is a single search hit.
type Result struct {
	RelPath string
	Title   string
	Excerpt string // surrounding context line
	Score   int    // higher = better match
}

// FullText searches every note's body for query string (case-insensitive).
// Returns results sorted by score (match count descending).
func FullText(query string) []Result {
	if query == "" {
		return nil
	}
	notesRoot := storage.NotesDir()
	q := strings.ToLower(query)
	var results []Result

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

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		content := string(data)
		lower := strings.ToLower(content)
		count := strings.Count(lower, q)
		if count == 0 {
			return nil
		}

		rel, _ := filepath.Rel(notesRoot, path)
		meta, _, _ := frontmatter.Parse(path)
		title := meta.Title
		if title == "" {
			title = strings.TrimSuffix(filepath.Base(path), ".md")
		}

		excerpt := extractExcerpt(content, query)

		results = append(results, Result{
			RelPath: rel,
			Title:   title,
			Excerpt: excerpt,
			Score:   count,
		})
		return nil
	})

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results
}

// Fuzzy filters a list of note paths by a fuzzy pattern against their name.
// Scores subsequence matches — typed chars must appear in order in the target.
func Fuzzy(pattern string, notes []string) []Result {
	if pattern == "" {
		out := make([]Result, 0, len(notes))
		for _, n := range notes {
			out = append(out, Result{RelPath: n, Score: 0})
		}
		return out
	}
	p := strings.ToLower(pattern)
	var results []Result
	for _, note := range notes {
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

// fuzzyScore returns a positive score if pattern is a subsequence of target.
// Consecutive matches and prefix matches score higher.
func fuzzyScore(pattern, target string) int {
	if pattern == "" {
		return 1
	}
	pi := 0
	score := 0
	consecutive := 0
	pr := []rune(pattern)
	tr := []rune(target)

	for ti, tc := range tr {
		if pi < len(pr) && unicode.ToLower(tc) == unicode.ToLower(pr[pi]) {
			pi++
			consecutive++
			score += consecutive * 2
			if ti == 0 {
				score += 5 // prefix bonus
			}
		} else {
			consecutive = 0
		}
	}
	if pi != len(pr) {
		return 0 // not a subsequence
	}
	// exact match bonus
	if target == pattern {
		score += 100
	}
	return score
}

// extractExcerpt finds the line containing query and returns a trimmed snippet.
func extractExcerpt(content, query string) string {
	lines := strings.Split(content, "\n")
	q := strings.ToLower(query)
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), q) {
			trimmed := strings.TrimSpace(line)
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
			// skip frontmatter lines
			if strings.HasPrefix(trimmed, "---") || strings.Contains(trimmed, ":") && len(trimmed) < 40 {
				continue
			}
			return trimmed
		}
	}
	return ""
}

// AllNotePaths returns relative paths of all .md files in the vault.
func AllNotePaths() []string {
	notesRoot := storage.NotesDir()
	var paths []string
	_ = filepath.Walk(notesRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.Contains(path, "/.git/") {
			return filepath.SkipDir
		}
		if strings.HasSuffix(path, ".md") {
			rel, _ := filepath.Rel(notesRoot, path)
			paths = append(paths, rel)
		}
		return nil
	})
	sort.Strings(paths)
	return paths
}
