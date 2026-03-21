package cmd

import (
	"fmt"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/config"
	"github.com/spf13/cobra"
)

var pinCmd = &cobra.Command{
	Use:   "pin <note>",
	Short: "Pin a note (marks it with ★ in the vault tree)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Pin(args[0]); err != nil {
			return err
		}
		fmt.Printf("★ pinned %s\n", args[0])
		return nil
	},
}

var unpinCmd = &cobra.Command{
	Use:   "unpin <note>",
	Short: "Unpin a note",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Unpin(args[0]); err != nil {
			return err
		}
		fmt.Printf("unpinned %s\n", args[0])
		return nil
	},
}

var pinnedCmd = &cobra.Command{
	Use:   "pinned",
	Short: "List all pinned notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.Load()
		if len(cfg.Pinned) == 0 {
			fmt.Println("No pinned notes. Use: gpad pin <note>")
			return nil
		}
		fmt.Println(colBold + "★ Pinned notes" + colReset)
		for _, p := range cfg.Pinned {
			fmt.Printf("  %s%s%s\n", colYellow, p, colReset)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pinCmd, unpinCmd, pinnedCmd)
}
