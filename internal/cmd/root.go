package cmd

import (
	"fmt"
	"os"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gpad",
	Short: "gpad - A lightning-fast CLI note manager",
	Long: `A modular, Git-synced markdown note manager designed 
for speed and terminal-centric workflows.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Ensure storage exists before any command runs
		return storage.EnsureDirs()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
}
