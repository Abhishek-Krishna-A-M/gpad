package cmd

import (
	"fmt"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/config"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/core"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Pull and push notes to the git remote",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.Load()
		if !cfg.GitEnabled {
			fmt.Println("Git sync not configured. Run: gpad git init <remote-url>")
			return nil
		}
		fmt.Println("Syncing...")
		return core.Sync()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
