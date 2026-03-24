package tui

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// ── ANSI escape sequences ─────────────────────────────────────────────────────

const (
	aReset    = "\033[0m"
	aBold     = "\033[1m"
	aDim      = "\033[2m"
	aItalic   = "\033[3m"
	aRev      = "\033[7m"
	aClearScr = "\033[2J"
	aHome     = "\033[H"
	aHideCur  = "\033[?25l"
	aShowCur  = "\033[?25h"
	aClearEOL = "\033[K"
	aClearDn  = "\033[J"

	// Tokyo Night — foreground colours only; no backgrounds so terminal
	// transparency (picom/compton on st) is fully respected.
	aFg     = "\033[38;2;192;202;245m" // #c0caf5  primary text
	aFgDim  = "\033[38;2;86;95;137m"   // #565f89  muted text
	aFgMut  = "\033[38;2;65;72;104m"   // #414868  very muted / divider
	aBlue   = "\033[38;2;122;162;247m" // #7aa2f7  dirs, headings
	aCyan   = "\033[38;2;115;218;202m" // #73daca  wikilinks, cyan accents
	aGreen  = "\033[38;2;158;206;106m" // #9ece6a  success, done items
	aYellow = "\033[38;2;224;175;104m" // #e0af68  tags, warnings, pinned
	aRed    = "\033[38;2;247;118;142m" // #f7768e  errors, delete
	aPurple = "\033[38;2;187;154;247m" // #bb9af7  code, blockquotes
	aOrange = "\033[38;2;255;158;100m" // #ff9e64  headings h2
)

// ── Cursor positioning ────────────────────────────────────────────────────────

// mc (moveCursor) returns the ANSI sequence to position the cursor
// at row, col (both 1-indexed). This is the ONLY function that should
// emit cursor positioning escapes — all draw code goes through here.
func mc(row, col int) string {
	return fmt.Sprintf("\033[%d;%dH", row, col)
}

// ── String measurement ────────────────────────────────────────────────────────

// vlen returns the visible rune width of s with ANSI codes stripped.
// Used everywhere widths are computed — special chars, box-drawing, CJK safe.
func vlen(s string) int {
	return len([]rune(stripANSI(s)))
}

// stripANSI removes all ANSI escape sequences from s.
func stripANSI(s string) string {
	var out strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == 0x1b && i+1 < len(s) && s[i+1] == '[' {
			i += 2
			for i < len(s) && s[i] != 'm' {
				i++
			}
			i++ // skip 'm'
			continue
		}
		out.WriteByte(s[i])
		i++
	}
	return out.String()
}

// truncANSI trims s to n visible runes, preserving ANSI escape sequences.
func truncANSI(s string, n int) string {
	if n <= 0 {
		return ""
	}
	var out strings.Builder
	visible := 0
	i := 0
	for i < len(s) && visible < n {
		if s[i] == 0x1b && i+1 < len(s) && s[i+1] == '[' {
			j := i + 2
			for j < len(s) && s[j] != 'm' {
				j++
			}
			out.WriteString(s[i : j+1])
			i = j + 1
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		out.WriteRune(r)
		visible++
		i += size
	}
	return out.String()
}

// pad returns s padded with spaces to exactly width visible runes.
// Truncates with … if s is longer than width.
func pad(s string, width int) string {
	vl := vlen(s)
	if vl == width {
		return s
	}
	if vl > width {
		if width <= 1 {
			return "…"
		}
		return truncANSI(s, width-1) + "…"
	}
	return s + strings.Repeat(" ", width-vl)
}

// rpt repeats s n times. Returns "" when n ≤ 0.
func rpt(s string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(s, n)
}
