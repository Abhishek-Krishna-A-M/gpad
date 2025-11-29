package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/config"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/gitrepo"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)

func handleDelete(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: gpad rm [-r] [--yes] <file_or_folder>")
		return
	}

	rec := false
	yes := false
	target := ""

	// parse flags
	for _, a := range args {
		switch a {
		case "-r", "--recursive":
			rec = true
		case "-y", "--yes":
			yes = true
		default:
			target = a
		}
	}

	if target == "" {
		fmt.Println("Error: no path provided")
		return
	}

	// Normalize and build full path inside notes dir
	notesRoot := storage.NotesDir()
	absTarget := filepath.Join(notesRoot, target)
	absTarget = filepath.Clean(absTarget)

	// Verify it's inside notes root
	if !strings.HasPrefix(absTarget, notesRoot) {
		fmt.Println("Error: deletion is only allowed inside", notesRoot)
		return
	}

	// Prevent deleting notes root itself
	if absTarget == notesRoot {
		fmt.Println("Error: refusing to delete the notes root directory")
		return
	}

	info, err := os.Stat(absTarget)
	if os.IsNotExist(err) {
		fmt.Println("Not found:", target)
		return
	}

	// If directory
	if info.IsDir() {
		handleDeleteDir(absTarget, target, rec, yes)
		return
	}

	// File deletion
	err = os.Remove(absTarget)
	if err != nil {
		fmt.Println("Failed to delete:", err)
		return
	}

	fmt.Println("Deleted:", target)
	commitAndPush("delete " + target)
}
func handleDeleteDir(absPath, rel string, rec bool, yes bool) {
	entries, _ := os.ReadDir(absPath)

	// not empty & no -r flag
	if len(entries) > 0 && !rec {
		fmt.Printf("Directory '%s' is not empty.\nUse: gpad rm -r %s\n", rel, rel)
		return
	}

	// recursive delete confirmation
	if rec && !yes {
		countFiles := 0
		countDirs := 0

		filepath.Walk(absPath, func(_ string, info os.FileInfo, _ error) error {
			if info.IsDir() {
				countDirs++
			} else {
				countFiles++
			}
			return nil
		})

		fmt.Printf("This will delete %d files and %d folders under '%s'.\n", countFiles, countDirs-1, rel)
		fmt.Print("Continue? (y/N): ")

		sc := bufio.NewScanner(os.Stdin)
		sc.Scan()
		resp := strings.ToLower(strings.TrimSpace(sc.Text()))

		if resp != "y" && resp != "yes" {
			fmt.Println("Aborted.")
			return
		}
	}

	err := os.RemoveAll(absPath)
	if err != nil {
		fmt.Println("Failed to delete directory:", err)
		return
	}

	fmt.Println("Deleted directory:", rel)
	commitAndPush("delete " + rel)
}

func commitAndPush(msg string) {
    cfg, _ := config.Load()
    if !cfg.AutoPush {
        return
    }

    notesRoot := storage.NotesDir()

    err := gitrepo.AddCommitPush(notesRoot, msg)
    if err != nil {
        fmt.Println("Git sync failed:", err)
    }
}


