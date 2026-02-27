package cmd

import (
	"fmt"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/gitrepo"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/config"
	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Manage note synchronization",
}

var initGitCmd = &cobra.Command{
	Use:   "init [remote-url]",
	Short: "Initialize git sync with a remote (SSH or HTTPS)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		path := storage.NotesDir()
		
		err := gitrepo.Initialize(path, url)
		if err != nil {
			return err
		}

		cfg, _ := config.Load()
		cfg.GitEnabled = true
		cfg.RepoURL = url
		cfg.AutoPush = true
		config.Save(cfg)

		fmt.Println("🚀 Git sync initialized and autopush enabled!")
		return nil
	},
}

func init() {
	gitCmd.AddCommand(initGitCmd)
	rootCmd.AddCommand(gitCmd)
}
