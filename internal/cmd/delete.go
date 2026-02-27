package cmd

import (
	"fmt"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/core"
	"github.com/spf13/cobra"
)

var (
	recursive bool
	force     bool
)

var deleteCmd = &cobra.Command{
	Use:     "rm [targets...]",
	Aliases: []string{"remove", "delete"},
	Short:   "Remove notes or directories",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Note: core.Delete should be updated to handle the 'force' logic 
		// or we handle the confirmation here.
		err := core.Delete(args, recursive)
		if err != nil {
			return err
		}

		fmt.Printf("Deleted %d item(s)\n", len(args))
		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "remove directories and their contents recursively")
	deleteCmd.Flags().BoolVarP(&force, "yes", "y", false, "skip confirmation prompt")
	rootCmd.AddCommand(deleteCmd)
}
