package cmd

import (
	"fmt"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/notes"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/ui"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [note]",
	Short: "Open a note (or search if no name provided)",
	RunE: func(cmd *cobra.Command, args []string) error {
		var target string

		if len(args) == 0 {
			// 1000x Speed: Logic for Fuzzy Finder
			allNotes, err := ui.GetAllNotes()
			if err != nil { return err }
			
			if len(allNotes) == 0 {
				fmt.Println("No notes found. Create one with 'gpad open my-note.md'")
				return nil
			}

			// For now, we print the list. 
			// In the next modular step, we'll pipe this to a fuzzy UI.
			fmt.Println("Select a note (or type the name):")
			for _, n := range allNotes {
				fmt.Printf(" - %s\n", n)
			}
			return nil
		}

		target = args[0]
		return notes.Open(target)
	},
	// This enables TAB completion in your shell!
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		list, _ := ui.GetAllNotes()
		return list, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
