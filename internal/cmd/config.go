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
var autopushCmd = &cobra.Command{
    Use:   "autopush [on/off]",
    Short: "Toggle automatic git pushing",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        val := args[0] == "on"
        cfg, _ := config.Load()
        cfg.AutoPush = val
        return config.Save(cfg)
    },
}
func init() {
	configCmd.AddCommand(setEditorCmd)
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(autopushCmd)
}
