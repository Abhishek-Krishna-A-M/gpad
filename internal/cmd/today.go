package cmd

import (
	"fmt"
	"time"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/core"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/daily"
	"github.com/spf13/cobra"
)

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Open today's daily note",
	Long: `Open (or create) today's daily note in notes/daily/YYYY-MM-DD.md.

Subcommands:
  gpad today           open today
  gpad today yesterday open yesterday
  gpad today list      show recent daily notes`,
}

var todayOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open today's daily note (default)",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = core.Sync()
		if err := daily.Open(); err != nil {
			return err
		}
		core.AutoSave("daily " + time.Now().Format("2006-01-02"))
		return nil
	},
}

var yesterdayCmd = &cobra.Command{
	Use:   "yesterday",
	Short: "Open yesterday's daily note",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = core.Sync()
		if err := daily.Yesterday(); err != nil {
			return err
		}
		core.AutoSave("daily yesterday")
		return nil
	},
}

var dailyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent daily notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		recent := daily.List(14)
		if len(recent) == 0 {
			fmt.Println("No daily notes yet.")
			return nil
		}
		fmt.Println("Recent daily notes:")
		for _, p := range recent {
			fmt.Printf("  %s\n", p)
		}
		return nil
	},
}

func init() {
	// gpad today → opens today directly (default run)
	todayCmd.RunE = func(cmd *cobra.Command, args []string) error {
		_ = core.Sync()
		if err := daily.Open(); err != nil {
			return err
		}
		core.AutoSave("daily " + time.Now().Format("2006-01-02"))
		return nil
	}
	todayCmd.AddCommand(todayOpenCmd, yesterdayCmd, dailyListCmd)
	rootCmd.AddCommand(todayCmd)
}
