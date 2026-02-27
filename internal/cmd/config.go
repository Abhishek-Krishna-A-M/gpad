package cmd

import (
	"fmt"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage gpad configuration",
}

var setEditorCmd = &cobra.Command{
	Use:   "editor [name]",
	Short: "Set default editor (vim, nvim, code, etc.)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.Load()
		cfg.Editor = args[0]
		err := config.Save(cfg)
		if err == nil {
			fmt.Printf("Editor set to: %s\n", args[0])
		}
		return err
	},
}

func init() {
	configCmd.AddCommand(setEditorCmd)
	rootCmd.AddCommand(configCmd)
}
