package cmd

import (
	"fmt"
	"os"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/core"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/help"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/viewer"
	"github.com/spf13/cobra"
)

// ── mv ──────────────────────────────────────────────────────────────────────

var moveCmd = &cobra.Command{
	Use:     "mv [sources...] <destination>",
	Aliases: []string{"move", "rename"},
	Short:   "Move or rename notes/folders",
	Long: `Move one or more notes to a destination, or rename a single note.
The H1 title in moved .md files is updated automatically.`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		dest := args[len(args)-1]
		srcs := args[:len(args)-1]
		if err := core.Move(srcs, dest); err != nil {
			return err
		}
		if len(srcs) == 1 {
			fmt.Printf("Moved %s → %s\n", srcs[0], dest)
		} else {
			fmt.Printf("Moved %d items → %s\n", len(srcs), dest)
		}
		return nil
	},
}

// ── cp ──────────────────────────────────────────────────────────────────────

var cpCmd = &cobra.Command{
	Use:   "cp <source> <destination>",
	Short: "Copy a note",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := core.Copy(args[0], args[1]); err != nil {
			return err
		}
		fmt.Printf("Copied %s → %s\n", args[0], args[1])
		return nil
	},
}

// ── rm ──────────────────────────────────────────────────────────────────────

var (
	rmRecursive bool
	rmForce     bool
)

var deleteCmd = &cobra.Command{
	Use:     "rm [targets...]",
	Aliases: []string{"remove", "delete"},
	Short:   "Remove notes or directories",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !rmForce {
			fmt.Printf("Delete %v? [y/N] ", args)
			var resp string
			fmt.Scanln(&resp)
			if resp != "y" && resp != "Y" {
				fmt.Println("Aborted.")
				return nil
			}
		}
		if err := core.Delete(args, rmRecursive); err != nil {
			return err
		}
		fmt.Printf("Deleted %d item(s)\n", len(args))
		return nil
	},
}

// ── help markdown ────────────────────────────────────────────────────────────

var helpMdCmd = &cobra.Command{
	Use:   "markdown",
	Short: "Show markdown & gpad syntax guide",
	Run: func(cmd *cobra.Command, args []string) {
		viewer.ViewRaw(help.GetMarkdownGuide())
	},
}

// ── completion ───────────────────────────────────────────────────────────────

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `Generate a shell completion script and source it to enable tab completion.
 
Bash (add to ~/.bashrc):
  source <(gpad completion bash)
 
Zsh (add to ~/.zshrc):
  source <(gpad completion zsh)
 
Fish:
  gpad completion fish | source
 
PowerShell:
  gpad completion powershell | Out-String | Invoke-Expression`,
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			// V2 is self-contained — no bash-completion package required
			rootCmd.GenBashCompletionV2(os.Stdout, true)
		case "zsh":
			rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&rmRecursive, "recursive", "r", false, "remove directories recursively")
	deleteCmd.Flags().BoolVarP(&rmForce, "yes", "y", false, "skip confirmation prompt")

	rootCmd.AddCommand(moveCmd, cpCmd, deleteCmd, completionCmd)
	rootCmd.AddCommand(helpMdCmd)
}
