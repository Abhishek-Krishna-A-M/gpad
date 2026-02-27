package cmd

import (
	"github.com/Abhishek-Krishna-A-M/gpad/internal/core"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync notes with the remote repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.Sync() // This calls the function in your uploaded sync.go
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
