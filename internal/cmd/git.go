package cmd

import (
	"fmt"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/config"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/gitrepo"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Git sync management",
}

var gitInitCmd = &cobra.Command{
	Use:   "init <remote-url>",
	Short: "Connect your vault to a git remote (SSH or HTTPS)",
	Long: `Connect the notes vault to a git remote for sync across machines.

Works with both SSH and HTTPS remotes:
  gpad git init git@github.com:user/notes.git
  gpad git init https://github.com/user/notes.git

gpad works fully offline without this step — git is optional.
Once connected, use 'gpad sync' to pull/push, or enable autopush
to have every save automatically synced.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		path := storage.NotesDir()

		fmt.Printf("Connecting vault to %s...\n", url)

		if err := gitrepo.Initialize(path, url); err != nil {
			return err
		}

		cfg, _ := config.Load()
		cfg.GitEnabled = true
		cfg.RepoURL = url
		cfg.AutoPush = true
		_ = config.Save(cfg)

		fmt.Println()
		fmt.Println("  Vault connected.")
		fmt.Println("  autopush is on — every save syncs automatically.")
		fmt.Println("  Run 'gpad sync' any time to manually pull/push.")
		fmt.Println()
		return nil
	},
}

var gitStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show git sync status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.Load()
		if !cfg.GitEnabled {
			fmt.Println("Git sync is not configured.")
			fmt.Println("Run: gpad git init <remote-url>")
			return nil
		}
		fmt.Printf("Remote: %s\n", cfg.RepoURL)
		fmt.Printf("Autopush: %v\n", cfg.AutoPush)
		return nil
	},
}

func init() {
	gitCmd.AddCommand(gitInitCmd, gitStatusCmd)
	rootCmd.AddCommand(gitCmd)
}
