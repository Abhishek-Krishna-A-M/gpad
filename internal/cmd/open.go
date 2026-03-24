package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/core"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/index"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/links"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/notes"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/tui"
	"github.com/spf13/cobra"
)

var wikilinkArg = regexp.MustCompile(`^\[\[(.+)\]\]$`)

var openCmd = &cobra.Command{
	Use:   "open [note]",
	Short: "Open or create a note — launches TUI when no argument given",
	Long: `Open a note by name or wikilink. With no argument, opens the full TUI.

  gpad open                    → launch full TUI
  gpad open my-note.md         → open by path
  gpad open '[[wikilink]]'     → resolve and open by wikilink
  gpad open ideas/quantum.md   → open in subfolder`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// No argument → launch the full TUI (same as bare gpad)
		if len(args) == 0 {
			return tui.Run()
		}

		target := args[0]

		// [[wikilink]] syntax — must be quoted in shell: gpad open '[[math]]'
		if m := wikilinkArg.FindStringSubmatch(target); m != nil {
			return openWikilink(m[1])
		}

		_ = core.Sync()
		if err := notes.Open(target); err != nil {
			return err
		}
		core.AutoSave("update " + target)
		return nil
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		all := index.AllNotes()
		paths := make([]string, 0, len(all))
		for _, n := range all {
			paths = append(paths, n.RelPath)
		}
		return paths, cobra.ShellCompDirectiveNoFileComp
	},
}

// openWikilink resolves [[inner]] syntax and opens the target note.
func openWikilink(inner string) error {
	target := inner
	if idx := strings.Index(target, "|"); idx != -1 {
		target = target[:idx]
	}
	if idx := strings.Index(target, "#"); idx != -1 {
		target = target[:idx]
	}
	target = strings.TrimSpace(target)

	absPath, found := links.Resolve(target)
	if !found {
		fmt.Printf("Note %q not found. Create it? [y/N] ", target)
		var resp string
		fmt.Scanln(&resp)
		if resp == "y" || resp == "Y" {
			name := strings.ToLower(strings.ReplaceAll(target, " ", "-")) + ".md"
			_ = core.Sync()
			if err := notes.Open(name); err != nil {
				return err
			}
			core.AutoSave("new " + name)
		}
		return nil
	}

	_ = core.Sync()
	if err := notes.Open(absPath); err != nil {
		return err
	}
	core.AutoSave("update " + absPath)
	return nil
}

func init() {
	rootCmd.AddCommand(openCmd)
}
