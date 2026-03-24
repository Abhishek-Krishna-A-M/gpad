package cmd

import (
	"fmt"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/search"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/tui"
	"github.com/spf13/cobra"
)

const (
	colReset  = "\033[0m"
	colBold   = "\033[1m"
	colCyan   = "\033[96m"
	colYellow = "\033[93m"
	colDim    = "\033[2m"
	colGreen  = "\033[92m"
)

var (
	fuzzyOnly bool
	ftsOnly   bool
)

var findCmd = &cobra.Command{
	Use:   "find [query]",
	Short: "Search notes — launches TUI search panel when no query given",
	Long: `Search across all notes.

  gpad find              → TUI with search panel open
  gpad find <query>      → full-text body search + fuzzy title match
  gpad find -f <query>   → fuzzy title only
  gpad find -t <query>   → full-text body only`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// No query → open TUI with search panel active
		if len(args) == 0 {
			return tui.Run()
		}

		query := strings.Join(args, " ")
		if !ftsOnly {
			runFuzzy(query)
		}
		if !fuzzyOnly {
			runFTS(query)
		}
		return nil
	},
}

func runFuzzy(query string) {
	all := search.AllNotePaths()
	results := search.Fuzzy(query, all)
	if len(results) == 0 {
		return
	}
	fmt.Printf("\n%s%s Title matches%s\n", colBold, colCyan, colReset)
	fmt.Println(colDim + strings.Repeat("─", 40) + colReset)
	for i, r := range results {
		if i >= 10 {
			fmt.Printf(colDim+"  … and %d more\n"+colReset, len(results)-10)
			break
		}
		fmt.Printf("  %s%s%s\n", colGreen, r.RelPath, colReset)
	}
}

func runFTS(query string) {
	results := search.FullText(query)
	if len(results) == 0 {
		if !fuzzyOnly {
			fmt.Printf(colDim+"\n  no body matches for %q\n"+colReset, query)
		}
		return
	}
	fmt.Printf("\n%s%s Body matches%s\n", colBold, colYellow, colReset)
	fmt.Println(colDim + strings.Repeat("─", 40) + colReset)
	for i, r := range results {
		if i >= 15 {
			fmt.Printf(colDim+"  … and %d more\n"+colReset, len(results)-15)
			break
		}
		fmt.Printf("  %s%s%s", colGreen, r.RelPath, colReset)
		if r.Title != "" && r.Title != r.RelPath {
			fmt.Printf("  %s%s%s", colDim, r.Title, colReset)
		}
		fmt.Printf("  ×%d\n", r.Score)
		if r.Excerpt != "" {
			fmt.Printf("    %s%s%s\n", colDim, r.Excerpt, colReset)
		}
	}
}

func init() {
	findCmd.Flags().BoolVarP(&fuzzyOnly, "fuzzy", "f", false, "fuzzy title search only")
	findCmd.Flags().BoolVarP(&ftsOnly, "text", "t", false, "full-text body search only")
	rootCmd.AddCommand(findCmd)
}
