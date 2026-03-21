package cmd

import (
	"fmt"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/frontmatter"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/tags"
	"github.com/spf13/cobra"
)

var tagsCmd = &cobra.Command{
	Use:   "tags [tag]",
	Short: "Browse tags across the vault",
	Long: `Without arguments: list every tag with note counts.
With a tag argument: list all notes carrying that tag.

Examples:
  gpad tags              → full tag index
  gpad tags programming  → notes tagged #programming`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return listAllTags()
		}
		return listNotesForTag(args[0])
	},
}

func listAllTags() error {
	all := tags.AllTags()
	if len(all) == 0 {
		fmt.Println(colDim + "  no tags yet — add tags: [] to frontmatter or use #inline in body" + colReset)
		return nil
	}
	idx := tags.Build()
	fmt.Printf("\n%s%s Tag index%s\n", colBold, colYellow, colReset)
	fmt.Println(colDim + strings.Repeat("─", 40) + colReset)
	for _, t := range all {
		notes := idx[t]
		fmt.Printf("  %s#%-24s%s %s%d%s\n",
			colYellow, t, colReset,
			colDim, len(notes), colReset)
	}
	fmt.Println()
	return nil
}

func listNotesForTag(tag string) error {
	tag = strings.TrimPrefix(tag, "#")
	notes := tags.NotesForTag(tag)
	if len(notes) == 0 {
		fmt.Printf(colDim+"  no notes tagged #%s\n"+colReset, tag)
		return nil
	}
	fmt.Printf("\n%s%s#%s%s\n", colBold, colYellow, tag, colReset)
	fmt.Println(colDim + strings.Repeat("─", 40) + colReset)
	for _, n := range notes {
		fmt.Printf("  %s%s%s\n", colGreen, n, colReset)
	}
	fmt.Println()
	return nil
}

// tag add / tag rm subcommands
var tagAddCmd = &cobra.Command{
	Use:   "add <tag> <note>",
	Short: "Add a tag to a note's frontmatter",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		tag, relPath := args[0], args[1]
		absPath := storage.AbsPath(relPath)
		if err := frontmatter.AddTag(absPath, tag); err != nil {
			return err
		}
		fmt.Printf("Added #%s to %s\n", strings.TrimPrefix(tag, "#"), relPath)
		return nil
	},
}

var tagRmCmd = &cobra.Command{
	Use:   "rm <tag> <note>",
	Short: "Remove a tag from a note's frontmatter",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		tag, relPath := args[0], args[1]
		absPath := storage.AbsPath(relPath)
		if err := frontmatter.RemoveTag(absPath, tag); err != nil {
			return err
		}
		fmt.Printf("Removed #%s from %s\n", strings.TrimPrefix(tag, "#"), relPath)
		return nil
	},
}

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Tag management (add/rm)",
}

func init() {
	tagCmd.AddCommand(tagAddCmd, tagRmCmd)
	rootCmd.AddCommand(tagsCmd, tagCmd)
}
