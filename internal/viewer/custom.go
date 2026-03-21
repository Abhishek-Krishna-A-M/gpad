package viewer

import (
	"fmt"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/links"
)

// ANSI palette
const (
	Bold    = "\033[1m"
	Dim     = "\033[2m"
	Cyan    = "\033[96m"
	Yellow  = "\033[93m"
	Magenta = "\033[95m"
	Green   = "\033[92m"
	Blue    = "\033[94m"
	Red     = "\033[91m"
	Reset   = "\033[0m"
)

// RenderCustom renders markdown to an ANSI-coloured string.
func RenderCustom(input string) (string, error) {
	lines := strings.Split(input, "\n")
	var out strings.Builder

	inCode := false
	inFrontmatter := false
	frontmatterDone := false
	emptyStreak := 0

	write := func(s string) {
		if strings.TrimSpace(s) == "" {
			emptyStreak++
			if emptyStreak > 1 {
				return
			}
		} else {
			emptyStreak = 0
		}
		out.WriteString(s + "\n")
	}

	for i, raw := range lines {
		line := strings.TrimSpace(raw)

		// --- frontmatter ---
		if i == 0 && line == "---" {
			inFrontmatter = true
			write(Dim + "──────────────── metadata ────────────────" + Reset)
			continue
		}
		if inFrontmatter {
			if line == "---" {
				inFrontmatter = false
				frontmatterDone = true
				write(Dim + "──────────────────────────────────────────" + Reset)
				continue
			}
			// render frontmatter key: value
			if idx := strings.Index(line, ":"); idx != -1 {
				key := line[:idx]
				val := strings.TrimSpace(line[idx+1:])
				write(fmt.Sprintf("%s%s%s: %s%s%s", Dim, key, Reset, Yellow, val, Reset))
			} else {
				write(Dim + line + Reset)
			}
			continue
		}
		_ = frontmatterDone

		// --- code fence ---
		if strings.HasPrefix(line, "```") {
			inCode = !inCode
			lang := strings.TrimPrefix(line, "```")
			if inCode {
				label := "code"
				if lang != "" {
					label = lang
				}
				dashes := 30 - len(label)
				if dashes < 0 {
					dashes = 0
				}
				write(Magenta + "┌─ " + label + " " + strings.Repeat("─", dashes) + Reset)
			} else {
				write(Magenta + "└" + strings.Repeat("─", 32) + Reset)
			}
			continue
		}
		if inCode {
			write(Magenta + raw + Reset)
			continue
		}

		// --- headings ---
		if strings.HasPrefix(line, "# ") {
			text := renderInline(strings.TrimSpace(line[2:]))
			write("\n" + Bold + Cyan + text + Reset)
			write(Cyan + strings.Repeat("═", len([]rune(line[2:]))+1) + Reset)
			continue
		}
		if strings.HasPrefix(line, "## ") {
			text := renderInline(strings.TrimSpace(line[3:]))
			write(Bold + Yellow + text + Reset)
			write(Yellow + strings.Repeat("─", len([]rune(line[3:]))+1) + Reset)
			continue
		}
		if strings.HasPrefix(line, "### ") {
			text := renderInline(strings.TrimSpace(line[4:]))
			write(Bold + Green + "▸ " + text + Reset)
			continue
		}
		if strings.HasPrefix(line, "#### ") {
			text := renderInline(strings.TrimSpace(line[5:]))
			write(Bold + Blue + "  ▸ " + text + Reset)
			continue
		}

		// --- task list ---
		if strings.HasPrefix(line, "- [ ] ") {
			write("  ○ " + renderInline(line[6:]))
			continue
		}
		if strings.HasPrefix(line, "- [x] ") || strings.HasPrefix(line, "- [X] ") {
			write(Dim + "  ✓ " + line[6:] + Reset)
			continue
		}

		// --- unordered list ---
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			write("  • " + renderInline(line[2:]))
			continue
		}
		// nested list
		if strings.HasPrefix(raw, "  - ") || strings.HasPrefix(raw, "  * ") {
			write("    ◦ " + renderInline(strings.TrimSpace(raw)[2:]))
			continue
		}

		// --- ordered list ---
		if len(line) > 2 && line[1] == '.' && line[0] >= '1' && line[0] <= '9' {
			write("  " + Blue + string(line[0]) + "." + Reset + " " + renderInline(line[3:]))
			continue
		}

		// --- blockquote ---
		if strings.HasPrefix(line, "> ") {
			write(Magenta + "┃ " + renderInline(line[2:]) + Reset)
			continue
		}

		// --- horizontal rule ---
		if line == "---" || line == "***" || line == "___" {
			write(Dim + strings.Repeat("─", 44) + Reset)
			continue
		}

		// --- bold / italic inline in paragraph ---
		if line != "" {
			write(renderInline(raw))
		} else {
			write("")
		}
	}

	return out.String(), nil
}

// renderInline processes **bold**, *italic*, `code`, [[wikilinks]], #tags inline.
func renderInline(s string) string {
	// [[wikilinks]]
	s = links.ReplaceForDisplay(s)

	// inline #tags
	words := strings.Fields(s)
	rendered := make([]string, 0, len(words))
	for _, w := range words {
		if strings.HasPrefix(w, "#") && len(w) > 1 {
			rendered = append(rendered, Yellow+"#"+strings.TrimPrefix(w, "#")+Reset)
		} else {
			rendered = append(rendered, w)
		}
	}
	s = strings.Join(rendered, " ")

	// **bold**
	s = applyDelim(s, "**", Bold, Reset)
	// *italic*
	s = applyDelim(s, "*", "\033[3m", Reset)
	// `code`
	s = applyDelim(s, "`", Magenta, Reset)

	return s
}

// applyDelim wraps text between paired delimiters with ANSI codes.
func applyDelim(s, delim, open, close string) string {
	parts := strings.Split(s, delim)
	if len(parts) < 3 {
		return s
	}
	var b strings.Builder
	for i, p := range parts {
		if i%2 == 1 {
			b.WriteString(open + p + close)
		} else {
			b.WriteString(p)
		}
	}
	return b.String()
}


