package cmd

import (
	"fmt"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/core"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/notes"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/ui"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [note]",
	Short: "Open or create a note",
	Long: `Open an existing note or create a new one.
If no argument is given, lists all notes interactively.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			all, err := ui.GetAllNotes()
			if err != nil {
				return err
			}
			if len(all) == 0 {
				fmt.Println("No notes yet. Try: gpad open my-first-note.md")
				return nil
			}
			fmt.Println("Notes (pass a name to open one):")
			for _, n := range all {
				fmt.Printf("  %s\n", n)
			}
			return nil
		}

		target := args[0]

		// pull latest before editing
		_ = core.Sync()

		if err := notes.Open(target); err != nil {
			return err
		}

		// push after editing
		core.AutoSave("update " + target)
		return nil
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		list, _ := ui.GetAllNotes()
		return list, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
