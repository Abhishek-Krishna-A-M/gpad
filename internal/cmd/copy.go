package cmd

import (
	"fmt"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/core"
	"github.com/spf13/cobra"
)

var cpCmd = &cobra.Command{
	Use:   "cp [source] [destination]",
	Short: "Copy a note",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := core.Copy(args[0], args[1])
		if err == nil {
			fmt.Printf("Copied %s to %s\n", args[0], args[1])
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(cpCmd)
}
