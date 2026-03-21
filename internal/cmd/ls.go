package cmd

import (
	"github.com/Abhishek-Krishna-A-M/gpad/internal/notes"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "Show vault tree (pinned notes marked ★)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return notes.List()
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
