package cmd

import (
	"fmt"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/core"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/notes"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/templates"
	"github.com/spf13/cobra"
)

var newTemplate string

var newCmd = &cobra.Command{
	Use:   "new <note>",
	Short: "Create a note from a template",
	Long: `Create a new note, optionally from a named template.

Examples:
  gpad new meeting-2026-03-22.md
  gpad new ideas/quantum-foam.md -t idea
  gpad new standup.md -t meeting`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]

		if newTemplate == "" {
			// show available templates and use default
			available := templates.List()
			if len(available) > 0 {
				fmt.Printf("Available templates: %v\n", available)
				fmt.Println("Tip: use -t <name> to pick one. Using default 'note' template.")
			}
			newTemplate = "note"
		}

		_ = core.Sync()

		if err := notes.Create(target, newTemplate); err != nil {
			return err
		}
		core.AutoSave("new " + target)
		return nil
	},
}

func init() {
	newCmd.Flags().StringVarP(&newTemplate, "template", "t", "", "template name to use")
	rootCmd.AddCommand(newCmd)
}
