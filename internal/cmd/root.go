package cmd

import (
	"fmt"
	"os"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/templates"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/tui"
	"github.com/spf13/cobra"
)

const version = "2.0.0"

var rootCmd = &cobra.Command{
	Use:   "gpad",
	Short: "gpad — a terminal-native knowledge vault",
	Long: `gpad 2.0 — fast, Git-synced, Obsidian-inspired notes for people who live in the terminal.

  gpad              launch full-screen TUI (default)
  gpad today        open today's daily note
  gpad open <note>  open or create a note
  gpad find         full-text + fuzzy search
  gpad sync         git sync`,
	Version: version,
	// bare "gpad" with no subcommand launches the TUI
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.Run()
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := storage.EnsureDirs(); err != nil {
			return err
		}
		templates.EnsureDefaults()
		return nil
	},
	// disable the default completion command Cobra adds
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
