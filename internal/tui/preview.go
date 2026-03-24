package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/frontmatter"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/links"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/notes"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/tags"
)

// getPreview returns rendered preview lines for the current selection.
// Cached by abs path — only re-rendered on cursor movement.
func (a *App) getPreview() []string {
	n := a.selected()
	if n == nil {
		return []string{aFgDim + "  nothing selected" + aReset}
	}
	if n.kind == kindDir {
		return a.dirPreview(n)
	}
	if n.absPath == a.previewCache {
		return a.previewLines
	}
	a.previewCache = n.absPath
	a.previewLines = renderPreview(n.absPath, a.treeHeight())
	return a.previewLines
}

func (a *App) dirPreview(n *treeNode) []string {
	var lines []string

	// display header: path relative to vault root
	header := n.relPath
	if header == "" {
		header = filepath.Base(n.absPath)
	}
	lines = append(lines, aBold+aBlue+header+"/"+aReset)
	lines = append(lines, aFgDim+strings.Repeat("─", 30)+aReset)

	// Read directory contents directly from disk — n.children is only
	// populated for the cwd root; all other dir nodes are shallow stubs.
	entries, err := os.ReadDir(n.absPath)
	if err != nil {
		lines = append(lines, aRed+"  cannot read directory"+aReset)
		return lines
	}

	var dirs, files []string
	for _, e := range entries {
		name := e.Name()
		if name == ".git" || strings.HasPrefix(name, ".") {
			continue
		}
		if e.IsDir() {
			dirs = append(dirs, name)
		} else if strings.HasSuffix(name, ".md") {
			files = append(files, name)
		}
	}

	count := 0
	for _, d := range dirs {
		lines = append(lines, "  "+aBlue+"▶ "+d+"/"+aReset)
		count++
		if count >= 24 {
			lines = append(lines, aFgDim+"  … and more"+aReset)
			break
		}
	}
	for _, f := range files {
		lines = append(lines, "  "+aFg+f+aReset)
		count++
		if count >= 32 {
			lines = append(lines, aFgDim+fmt.Sprintf("  … %d more", len(dirs)+len(files)-count)+aReset)
			break
		}
	}

	if count == 0 {
		lines = append(lines, aFgDim+"  empty directory"+aReset)
	}

	// summary line
	lines = append(lines, "")
	lines = append(lines, aFgDim+fmt.Sprintf("  %d dirs · %d notes", len(dirs), len(files))+aReset)
	return lines
}

// renderPreview renders an .md file to ANSI-coloured lines for the right panel.
// Uses frontmatter.Parse() for metadata and notes.Stats() for word/line/link counts.
func renderPreview(absPath string, maxLines int) []string {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return []string{aRed + "  cannot read file" + aReset}
	}

	// Use frontmatter.Parse to properly split metadata from body
	meta, body, _ := frontmatter.Parse(absPath)

	var lines []string

	// render frontmatter as styled metadata block
	if meta.Title != "" || len(meta.Tags) > 0 || !meta.Date.IsZero() {
		lines = append(lines, aFgDim+"── meta ──"+aReset)
		if meta.Title != "" {
			lines = append(lines, aFgDim+"title"+aReset+": "+aYellow+meta.Title+aReset)
		}
		if !meta.Date.IsZero() {
			lines = append(lines, aFgDim+"date"+aReset+":  "+aFgDim+meta.Date.Format("2006-01-02")+aReset)
		}
		if len(meta.Tags) > 0 {
			ts := make([]string, len(meta.Tags))
			for i, t := range meta.Tags {
				ts[i] = aYellow+"#"+t+aReset
			}
			lines = append(lines, aFgDim+"tags"+aReset+":  "+strings.Join(ts, " "))
		}
		lines = append(lines, aFgDim+"──────────"+aReset)
		lines = append(lines, "")
	}

	// render body
	inCode := false
	emptyStreak := 0
	rawLines := strings.Split(body, "\n")

	for _, line := range rawLines {
		if len(lines) >= maxLines-6 {
			lines = append(lines, aFgDim+"  … (v to page)"+aReset)
			break
		}
		tr := strings.TrimSpace(line)

		if strings.HasPrefix(tr, "```") {
			inCode = !inCode
			lang := strings.TrimPrefix(tr, "```")
			if inCode {
				label := lang
				if label == "" { label = "code" }
				lines = append(lines, aPurple+"┌─ "+label+aReset)
			} else {
				lines = append(lines, aPurple+"└──────────"+aReset)
			}
			continue
		}
		if inCode {
			lines = append(lines, aPurple+line+aReset)
			continue
		}

		if tr == "" {
			emptyStreak++
			if emptyStreak > 1 { continue }
			lines = append(lines, "")
			continue
		}
		emptyStreak = 0

		switch {
		case strings.HasPrefix(tr, "# "):
			lines = append(lines, aBold+aBlue+strings.TrimSpace(tr[2:])+aReset)
		case strings.HasPrefix(tr, "## "):
			lines = append(lines, aBold+aOrange+strings.TrimSpace(tr[3:])+aReset)
		case strings.HasPrefix(tr, "### "):
			lines = append(lines, aBold+aGreen+"▸ "+strings.TrimSpace(tr[4:])+aReset)
		case strings.HasPrefix(tr, "#### "):
			lines = append(lines, aFgDim+"  ▸ "+strings.TrimSpace(tr[5:])+aReset)
		case strings.HasPrefix(tr, "- [ ] "):
			lines = append(lines, "  "+aFgDim+"○"+aReset+" "+renderInline(tr[6:]))
		case strings.HasPrefix(tr, "- [x] "), strings.HasPrefix(tr, "- [X] "):
			lines = append(lines, "  "+aGreen+"✓"+aReset+" "+aFgDim+tr[6:]+aReset)
		case strings.HasPrefix(tr, "- "), strings.HasPrefix(tr, "* "):
			lines = append(lines, "  • "+renderInline(tr[2:]))
		case strings.HasPrefix(tr, "> "):
			lines = append(lines, aPurple+"┃ "+renderInline(tr[2:])+aReset)
		case tr == "---" || tr == "***":
			lines = append(lines, aFgDim+strings.Repeat("─", 28)+aReset)
		default:
			lines = append(lines, renderInline(line))
		}
	}

	// footer: word count + links + tags from notes.Stats() and tags package
	lines = append(lines, "")
	lines = append(lines, aFgDim+strings.Repeat("─", 28)+aReset)

	relPath, _ := filepath.Rel(storage.NotesDir(), absPath)
	wordCount, lineCount, linkCount, err := notes.Stats(absPath)
	if err == nil {
		lines = append(lines, fmt.Sprintf("%s%d words · %d lines · %d links%s",
			aFgDim, wordCount, lineCount, linkCount, aReset))
	} else {
		lines = append(lines, fmt.Sprintf("%s%d words%s",
			aFgDim, len(strings.Fields(string(data))), aReset))
	}

	// tags from our tags package (frontmatter + inline)
	noteTags := tags.TagsForNote(relPath)
	if len(noteTags) > 0 {
		ts := make([]string, len(noteTags))
		for i, t := range noteTags {
			ts[i] = aYellow+"#"+t+aReset
		}
		lines = append(lines, strings.Join(ts, " "))
	}

	// backlinks from our links package
	back := links.Backlinks(relPath)
	if len(back) > 0 {
		shown := back
		if len(shown) > 3 { shown = shown[:3] }
		lines = append(lines, aGreen+"← "+strings.Join(shown, ", ")+aReset)
	}
	return lines
}

// renderInline handles **bold**, *italic*, `code`, [[wikilinks]], #tags.
func renderInline(s string) string {
	// [[wikilinks]] — use links package display format
	for strings.Contains(s, "[[") {
		st := strings.Index(s, "[[")
		en := strings.Index(s[st:], "]]")
		if en < 0 { break }
		en += st
		inner := s[st+2 : en]
		if idx := strings.Index(inner, "|"); idx != -1 {
			inner = inner[idx+1:]
		}
		if idx := strings.Index(inner, "#"); idx != -1 {
			inner = inner[:idx]
		}
		s = s[:st] + aCyan+"→ "+strings.TrimSpace(inner)+aReset + s[en+2:]
	}
	// inline #tags
	words := strings.Fields(s)
	for i, w := range words {
		if strings.HasPrefix(w, "#") && len(w) > 1 {
			words[i] = aYellow+w+aReset
		}
	}
	s = strings.Join(words, " ")
	s = applyDelim(s, "**", aBold, aReset)
	s = applyDelim(s, "*", aItalic, aReset)
	s = applyDelim(s, "`", aPurple, aReset)
	return s
}

func applyDelim(s, delim, open, close string) string {
	parts := strings.Split(s, delim)
	if len(parts) < 3 { return s }
	var b strings.Builder
	for i, p := range parts {
		if i%2 == 1 { b.WriteString(open + p + close) } else { b.WriteString(p) }
	}
	return b.String()
}
