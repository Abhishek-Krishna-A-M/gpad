package cmd

import (
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/viewer"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view [note]",
	Short: "View a note with markdown rendering",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		absPath := storage.AbsPath(args[0])
		return viewer.View(absPath)
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
