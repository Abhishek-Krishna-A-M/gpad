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

var configEditorCmd = &cobra.Command{
	Use:   "editor <name>",
	Short: "Set preferred editor (nvim, code, micro, ...)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.Load()
		cfg.Editor = args[0]
		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Printf("Editor set to: %s\n", args[0])
		return nil
	},
}

var configAutopushCmd = &cobra.Command{
	Use:   "autopush [on|off]",
	Short: "Toggle automatic git push after every save",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		val := args[0] == "on"
		cfg, _ := config.Load()
		cfg.AutoPush = val
		if err := config.Save(cfg); err != nil {
			return err
		}
		state := "off"
		if val {
			state = "on"
		}
		fmt.Printf("Autopush: %s\n", state)
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.Load()
		fmt.Printf("editor:    %s\n", cfg.Editor)
		fmt.Printf("git:       %v\n", cfg.GitEnabled)
		fmt.Printf("remote:    %s\n", cfg.RepoURL)
		fmt.Printf("autopush:  %v\n", cfg.AutoPush)
		fmt.Printf("pinned:    %v\n", cfg.Pinned)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configEditorCmd, configAutopushCmd, configShowCmd)
	rootCmd.AddCommand(configCmd)
}
