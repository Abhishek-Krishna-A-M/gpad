package cmd

import (
	"github.com/Abhishek-Krishna-A-M/gpad/internal/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:     "tui",
	Aliases: []string{"ui"},
	Short:   "Open the full-screen TUI (also the default when no command given)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.Run()
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
