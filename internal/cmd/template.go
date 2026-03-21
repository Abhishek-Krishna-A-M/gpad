package cmd

import (
	"fmt"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/editor"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/templates"
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage note templates",
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		list := templates.List()
		if len(list) == 0 {
			fmt.Println("No templates. They live in ~/.gpad/templates/")
			return nil
		}
		fmt.Println(colBold + "Templates" + colReset)
		for _, t := range list {
			fmt.Printf("  %s%s%s\n", colCyan, t, colReset)
		}
		return nil
	},
}

var templateNewCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create a new template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		starter := `---
title: {{title}}
date: {{date}}
tags: []
---

# {{title}}

{{cursor}}
`
		if err := templates.Save(name, starter); err != nil {
			return err
		}
		fmt.Printf("Created template %q — opening in editor...\n", name)
		return editor.Open(templates.Path(name))
	},
}

var templateEditCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "Edit an existing template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return editor.Open(templates.Path(args[0]))
	},
}

var templateDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := templates.Delete(args[0]); err != nil {
			return err
		}
		fmt.Printf("Deleted template %q\n", args[0])
		return nil
	},
}

func init() {
	templateCmd.AddCommand(templateListCmd, templateNewCmd, templateEditCmd, templateDeleteCmd)
	rootCmd.AddCommand(templateCmd)
}
