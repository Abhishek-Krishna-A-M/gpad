package links_test

import (
	"testing"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/links"
)

func TestParseTargets(t *testing.T) {
	cases := []struct {
		input string
		want  []string
	}{
		{
			input: "See [[math]] for details.",
			want:  []string{"math"},
		},
		{
			input: "Links to [[college/physics]] and [[ideas|brainstorm]].",
			want:  []string{"college/physics", "ideas"},
		},
		{
			input: "Anchor link [[note#heading]] here.",
			want:  []string{"note"},
		},
		{
			input: "No links here.",
			want:  nil,
		},
		{
			input: "Multiple [[a]], [[b]], [[c|alias]].",
			want:  []string{"a", "b", "c"},
		},
	}

	for _, tc := range cases {
		got := links.ParseTargets(tc.input)
		if len(got) != len(tc.want) {
			t.Errorf("input %q: got %v want %v", tc.input, got, tc.want)
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("input %q: index %d got %q want %q", tc.input, i, got[i], tc.want[i])
			}
		}
	}
}

func TestReplaceForDisplay(t *testing.T) {
	input := "See [[math]] and [[physics|Physics course]]."
	out := links.ReplaceForDisplay(input)

	if !containsStr(out, "→ math") {
		t.Errorf("expected '→ math' in output, got: %q", out)
	}
	if !containsStr(out, "→ Physics course") {
		t.Errorf("expected '→ Physics course' in output, got: %q", out)
	}
	if containsStr(out, "[[") {
		t.Errorf("raw [[ should not appear in output: %q", out)
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && func() bool {
		for i := 0; i <= len(s)-len(sub); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	}()
}
