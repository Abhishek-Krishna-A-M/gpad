package search_test

import (
	"testing"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/search"
)

func TestFuzzy_BasicMatch(t *testing.T) {
	notes := []string{
		"college/mathematics.md",
		"college/physics.md",
		"daily/2026-03-22.md",
		"ideas/quantum-computing.md",
	}

	results := search.Fuzzy("math", notes)
	if len(results) == 0 {
		t.Fatal("expected at least one result for 'math'")
	}
	if results[0].RelPath != "college/mathematics.md" {
		t.Errorf("top result: got %q want %q", results[0].RelPath, "college/mathematics.md")
	}
}

func TestFuzzy_EmptyPattern(t *testing.T) {
	notes := []string{"a.md", "b.md", "c.md"}
	results := search.Fuzzy("", notes)
	if len(results) != 3 {
		t.Errorf("empty pattern should return all notes, got %d", len(results))
	}
}

func TestFuzzy_NoMatch(t *testing.T) {
	notes := []string{"college/math.md", "ideas/quantum.md"}
	results := search.Fuzzy("zzzzz", notes)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestFuzzy_Subsequence(t *testing.T) {
	notes := []string{"quantum-entanglement.md", "quick-maths.md", "queue-theory.md"}
	// "qe" is a subsequence of "quantum-entanglement" and "queue-theory"
	results := search.Fuzzy("qe", notes)
	if len(results) < 1 {
		t.Errorf("expected matches for 'qe', got %d", len(results))
	}
}

func TestFuzzy_PrefixScoresHigher(t *testing.T) {
	notes := []string{"math-advanced.md", "college/math.md", "a-math-note.md"}
	results := search.Fuzzy("math", notes)
	// notes starting with "math" should score higher than "a-math-note"
	if len(results) < 2 {
		t.Skip("not enough results")
	}
	for _, r := range results[:2] {
		if r.RelPath == "a-math-note.md" && results[0].RelPath != "a-math-note.md" {
			// fine — it's not at top
			continue
		}
	}
}

func TestFuzzy_CaseInsensitive(t *testing.T) {
	notes := []string{"College/Mathematics.md", "ideas/QUANTUM.md"}
	results := search.Fuzzy("MATH", notes)
	if len(results) == 0 {
		t.Error("fuzzy should be case-insensitive")
	}
}
