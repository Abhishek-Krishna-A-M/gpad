package cmd

import (
	"fmt"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/core"
	"github.com/spf13/cobra"
)

var moveCmd = &cobra.Command{
	Use:     "mv [sources...] [destination]",
	Aliases: []string{"move", "rename"},
	Short:   "Move or rename notes/folders",
	Long: `Move one or more notes to a destination directory, or rename a single note.
If a note is moved or renamed, gpad automatically updates the Markdown H1 header.`,
	Args:    cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// The last argument is the destination
		destination := args[len(args)-1]
		// Everything before the last argument are the sources
		sources := args[:len(args)-1]

		err := core.Move(sources, destination)
		if err != nil {
			return err
		}

		if len(sources) > 1 {
			fmt.Printf("Successfully moved %d items to %s\n", len(sources), destination)
		} else {
			fmt.Printf("Successfully moved %s to %s\n", sources[0], destination)
		}
		
		return nil
	},
}

func init() {
	// Register the command to the root
	rootCmd.AddCommand(moveCmd)
}
