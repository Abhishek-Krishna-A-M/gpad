package cmd

import (
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/viewer"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view <note>",
	Short: "Render a note in the terminal",
	Long: `Render a note with markdown formatting, wikilinks, tags,
backlinks panel, and word-count stats. Uses less for scrolling.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return viewer.View(storage.AbsPath(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
