package cmd

import (
	"fmt"
	"os"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/templates"
	"github.com/spf13/cobra"
)

const version = "2.0.0"

var rootCmd = &cobra.Command{
	Use:   "gpad",
	Short: "gpad — a terminal-native knowledge vault",
	Long: `gpad 2.0 — fast, Git-synced, Obsidian-inspired notes for people who live in the terminal.

  gpad today              open today's daily note
  gpad open <note>        open or create a note
  gpad find <query>       full-text + fuzzy search
  gpad links <note>       backlinks & outlinks
  gpad tags               tag index
  gpad graph              ASCII link graph
  gpad ls                 vault tree`,
	Version: version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := storage.EnsureDirs(); err != nil {
			return err
		}
		templates.EnsureDefaults()
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
