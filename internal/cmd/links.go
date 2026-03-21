package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/links"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
	"github.com/spf13/cobra"
)

var linksCmd = &cobra.Command{
	Use:     "links [note]",
	Aliases: []string{"backlinks", "bl"},
	Short:   "Show backlinks and outlinks for a note",
	Long: `Display the full link graph for a note:
  outlinks  — notes this note links TO via [[wikilinks]]
  backlinks — notes that link TO this note`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		relPath := args[0]
		// normalise: strip leading notes/ prefix if user typed it
		notesRoot := storage.NotesDir()
		if filepath.IsAbs(relPath) {
			var err error
			relPath, err = filepath.Rel(notesRoot, relPath)
			if err != nil {
				return err
			}
		}

		out := links.Outlinks(relPath)
		back := links.Backlinks(relPath)

		fmt.Printf("\n%s%s%s\n", colBold, relPath, colReset)
		fmt.Println(colDim + strings.Repeat("─", 44) + colReset)

		if len(out) == 0 && len(back) == 0 {
			fmt.Println(colDim + "  no links found — add [[note name]] to connect notes" + colReset)
			return nil
		}

		if len(out) > 0 {
			fmt.Printf("%s  outlinks%s\n", colCyan, colReset)
			for _, o := range out {
				fmt.Printf("    %s→%s %s\n", colCyan, colReset, o)
			}
		}

		if len(back) > 0 {
			fmt.Printf("%s  backlinks%s\n", colGreen, colReset)
			for _, b := range back {
				fmt.Printf("    %s←%s %s\n", colGreen, colReset, b)
			}
		}
		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(linksCmd)
}
