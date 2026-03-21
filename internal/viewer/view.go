package viewer

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/editor"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/links"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/notes"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/tags"
)

// View renders a note with backlinks and stats, then pages it with less.
func View(absPath string) error {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	rendered, _ := RenderCustom(string(data))

	// Append backlinks panel
	rel, _ := filepath.Rel(storage.NotesDir(), absPath)
	rendered += buildBacklinksPanel(rel)
	rendered += buildStatsPanel(absPath)

	// Page through less -R (keeps ANSI colours)
	cmd := exec.Command("less", "-R", "--quit-if-one-screen")
	cmd.Stdin = strings.NewReader(rendered)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	return promptEdit(absPath)
}

// ViewRaw renders and prints without paging (for short output like help text).
func ViewRaw(text string) {
	if rendered, err := RenderCustom(text); err == nil {
		fmt.Print(rendered)
		return
	}
	fmt.Print(text)
}

func buildBacklinksPanel(relPath string) string {
	bl := links.Backlinks(relPath)
	ol := links.Outlinks(relPath)

	var b strings.Builder
	b.WriteString("\n" + Dim + strings.Repeat("─", 44) + Reset + "\n")
	b.WriteString(Bold + "  Links\n" + Reset)

	if len(ol) == 0 && len(bl) == 0 {
		b.WriteString(Dim + "  no links — add [[note name]] to connect notes\n" + Reset)
		return b.String()
	}

	if len(ol) > 0 {
		b.WriteString(Dim + "  outlinks\n" + Reset)
		for _, o := range ol {
			b.WriteString(fmt.Sprintf("    %s→%s %s\n", Cyan, Reset, o))
		}
	}
	if len(bl) > 0 {
		b.WriteString(Dim + "  backlinks\n" + Reset)
		for _, bk := range bl {
			b.WriteString(fmt.Sprintf("    %s←%s %s\n", Green, Reset, bk))
		}
	}
	return b.String()
}

func buildStatsPanel(absPath string) string {
	words, lines, linkCount, err := notes.Stats(absPath)
	if err != nil {
		return ""
	}

	rel, _ := filepath.Rel(storage.NotesDir(), absPath)
	noteTags := tags.TagsForNote(rel)

	var b strings.Builder
	b.WriteString(Dim + strings.Repeat("─", 44) + Reset + "\n")
	b.WriteString(fmt.Sprintf(Dim+"  %d words · %d lines · %d links"+Reset+"\n", words, lines, linkCount))
	if len(noteTags) > 0 {
		tagStr := ""
		for _, t := range noteTags {
			tagStr += Yellow + "#" + t + Reset + " "
		}
		b.WriteString("  " + strings.TrimSpace(tagStr) + "\n")
	}
	return b.String()
}

func promptEdit(path string) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n  edit? [y/N] ")
	resp, _ := reader.ReadString('\n')
	resp = strings.ToLower(strings.TrimSpace(resp))
	if resp == "y" || resp == "yes" {
		return editor.Open(path)
	}
	return nil
}
