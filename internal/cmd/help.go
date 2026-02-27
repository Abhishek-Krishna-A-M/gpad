package cmd

import (
	"github.com/Abhishek-Krishna-A-M/gpad/internal/help"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/viewer"
	"github.com/spf13/cobra"
)

var helpMdCmd = &cobra.Command{
	Use:   "markdown",
	Short: "Show Markdown syntax guide",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the raw help text
		rawHelp := help.GetMarkdownGuide() 
		
		// Use your viewer to make it look professional
		viewer.ViewRaw(rawHelp)
	},
}

func init() {
	// You can nest this under a 'help' subcommand or a root 'md' command
	rootCmd.AddCommand(helpMdCmd)
}
