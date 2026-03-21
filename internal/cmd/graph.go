package cmd

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/links"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/search"
	"github.com/spf13/cobra"
)

var graphCmd = &cobra.Command{
	Use:   "graph [note]",
	Short: "ASCII link graph of the vault",
	Long: `Visualise the wikilink graph in the terminal.

Without arguments: full vault graph.
With a note argument: ego graph (note + its immediate neighbours).

Examples:
  gpad graph
  gpad graph college/math.md`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fullGraph()
		}
		return egoGraph(args[0])
	},
}

func fullGraph() error {
	idx := links.BuildIndex()
	if len(idx) == 0 {
		fmt.Println(colDim + "  vault has no [[wikilinks]] yet" + colReset)
		return nil
	}

	// Sort notes for stable output
	notes := search.AllNotePaths()

	fmt.Printf("\n%s%s Link graph%s\n", colBold, colCyan, colReset)
	fmt.Println(colDim + strings.Repeat("─", 50) + colReset)

	// find max note name length for alignment
	maxLen := 0
	for _, n := range notes {
		if len(n) > maxLen {
			maxLen = len(n)
		}
	}
	if maxLen > 36 {
		maxLen = 36
	}

	// only show notes that participate in the graph
	hasLinks := map[string]bool{}
	for _, n := range notes {
		e, ok := idx[n]
		if ok && (len(e.OutLinks) > 0 || len(e.InLinks) > 0) {
			hasLinks[n] = true
		}
	}

	if len(hasLinks) == 0 {
		fmt.Println(colDim + "  no linked notes — add [[note name]] to connect notes" + colReset)
		return nil
	}

	for _, n := range notes {
		if !hasLinks[n] {
			continue
		}
		e := idx[n]
		label := truncate(n, maxLen)
		pad := strings.Repeat(" ", maxLen-len(label))

		outStr := ""
		if len(e.OutLinks) > 0 {
			targets := make([]string, len(e.OutLinks))
			copy(targets, e.OutLinks)
			sort.Strings(targets)
			outStr = colCyan + " → " + colReset + colDim + strings.Join(targets, ", ") + colReset
		}

		inCount := ""
		if len(e.InLinks) > 0 {
			inCount = fmt.Sprintf(colGreen+" ←%d"+colReset, len(e.InLinks))
		}

		fmt.Printf("  %s%s%s%s%s%s\n",
			colGreen, label, colReset,
			pad,
			inCount,
			outStr,
		)
	}
	fmt.Println()
	return nil
}

func egoGraph(relPath string) error {
	idx := links.BuildIndex()
	e, ok := idx[relPath]
	if !ok {
		fmt.Printf(colDim+"  %s not found in link graph\n"+colReset, relPath)
		return nil
	}

	name := filepath.Base(relPath)
	fmt.Printf("\n%s%s%s\n", colBold, name, colReset)
	fmt.Println(colDim + strings.Repeat("─", 40) + colReset)

	// Centre node
	fmt.Printf("  %s◉ %s%s\n", colCyan, relPath, colReset)

	if len(e.OutLinks) > 0 {
		fmt.Printf("\n  %slinks to%s\n", colCyan, colReset)
		sorted := sorted(e.OutLinks)
		for i, o := range sorted {
			branch := "├──"
			if i == len(sorted)-1 {
				branch = "└──"
			}
			fmt.Printf("    %s%s %s%s\n", colCyan, branch, colReset, o)
		}
	}

	if len(e.InLinks) > 0 {
		fmt.Printf("\n  %slinked from%s\n", colGreen, colReset)
		sorted := sorted(e.InLinks)
		for i, b := range sorted {
			branch := "├──"
			if i == len(sorted)-1 {
				branch = "└──"
			}
			fmt.Printf("    %s%s %s%s\n", colGreen, branch, colReset, b)
		}
	}
	fmt.Println()
	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func sorted(s []string) []string {
	c := make([]string, len(s))
	copy(c, s)
	sort.Strings(c)
	return c
}

func init() {
	rootCmd.AddCommand(graphCmd)
}
