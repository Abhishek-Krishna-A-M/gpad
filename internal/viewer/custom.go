package viewer

import (
	"strings"
)

// ANSI
const (
	Bold    = "\033[1m"
	Cyan    = "\033[96m"
	Yellow  = "\033[93m"
	Magenta = "\033[95m"
	Green   = "\033[92m"
	Reset   = "\033[0m"
)

func underline(s string, char string) string {
	return strings.Repeat(char, len(s))
}

func RenderCustom(input string) (string, error) {
	lines := strings.Split(input, "\n")
	var out strings.Builder

	inCode := false
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

	for _, raw := range lines {
		line := strings.TrimSpace(raw)

		if strings.HasPrefix(line, "```") {
			inCode = !inCode
			if inCode {
				write(Magenta + "┌─ code ───────────────────────" + Reset)
			} else {
				write(Magenta + "└──────────────────────────────" + Reset)
			}
			continue
		}

		if inCode {
			write(Magenta + raw + Reset)
			continue
		}

		// H1
		if strings.HasPrefix(line, "# ") {
			text := (strings.TrimSpace(line[2:])
			write(Bold + Cyan + text + Reset)
			continue
		}

		// H2
		if strings.HasPrefix(line, "## ") {
			text := strings.TrimSpace(line[3:])
			write(Bold + Yellow + text + Reset)
			continue
		}

		// H3
		if strings.HasPrefix(line, "### ") {
			text := strings.TrimSpace(line[4:])
			write(Bold + Green + "▸ " + text + Reset)
			continue
		}

		// List
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			write("• " + line[2:])
			continue
		}

		// Blockquote
		if strings.HasPrefix(line, "> ") {
			write(Magenta + "┃ " + line[2:] + Reset)
			continue
		}

		// Horizontal rule
		if line == "---" {
			write(Magenta + "──────────────────────────────" + Reset)
			continue
		}

		// Normal paragraph or blank
		write(raw)
	}

	return out.String(), nil
}
