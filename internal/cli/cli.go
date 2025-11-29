package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/notes"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/viewer"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/gitrepo"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/config"

)

func Run() {
	if err := storage.EnsureDirs(); err != nil {
		fmt.Println("Error:", err)
		return
	}
	args := os.Args[1:]

	if len(args) == 0 {
		printHelp()
		return
	}

	if len(args) == 1 && filepath.Ext(args[0]) == ".md" {
		handleOpen(args[0])
		return
	}

	switch args[0] {

	case "view":
		if len(args) < 2 {
			fmt.Println("Usage: gpad view <file>")
			return
		}
		handleView(args[1])

	case "list":
		handleList()

	case "init":
		handleInit(args[1:])

	case "config":
		handleConfig(args[1:])

	case "uninstall":
		handleUninstall(args[1:])

	default:
		fmt.Println("Unknown command:", args[0])
		printHelp()
	}
}

func printHelp() {
	fmt.Println(`gpad - global markdown notes

Usage:
  gpad <path.md>             Create or edit a note instantly
  gpad view <path.md>        View a markdown file (anywhere)
  gpad list                  List all notes
  gpad init [--github URL]   Initialize notes storage
  gpad config ...            Configure editor/autopush
  gpad uninstall             Remove gpad from system
`)
}

func handleOpen(path string) {
	if err := notes.Open(path); err != nil {
		fmt.Println("Error:", err)
	}
}
func handleView(path string) {
	if err := viewer.View(path); err != nil {
		fmt.Println("Error:", err)
	}
}
func handleList() {
	if err := notes.List(); err != nil {
		fmt.Println("Error:", err)
	}
}
func handleInit(args []string) {
	var repo string

	// Parse flag
	if len(args) >= 2 && args[0] == "--github" {
		repo = args[1]
	}

	notesPath := storage.NotesDir()

	// Offline init
	if repo == "" {
		fmt.Println("Initializing gpad in offline mode...")
		fmt.Println("Notes stored at:", notesPath)
		return
	}

	// Git mode
	fmt.Println("Initializing gpad with GitHub repo:", repo)

	if !gitrepo.Exists() {
		fmt.Println("Error: git is not installed.")
		return
	}

	tmpDir := filepath.Join(storage.GpadDir(), "tmp_clone")

	// Clean temp dir
	os.RemoveAll(tmpDir)
	os.MkdirAll(tmpDir, 0755)

	// Clone into tmp
	if err := gitrepo.Clone(repo, tmpDir); err != nil {
		fmt.Println("Git clone failed:", err)
		return
	}

	// Merge offline notes into tmp clone
	if dirNotEmpty(notesPath) {
		fmt.Println("Merging offline notes into GitHub repo...")
		if err := gitrepo.MergeOfflineIntoRepo(tmpDir, notesPath); err != nil {
			fmt.Println("Merge failed:", err)
			return
		}
	}

	// Replace notes folder with merged repo
	os.RemoveAll(notesPath)
	os.Rename(tmpDir, notesPath)

	// Auto commit merge
	fmt.Println("Syncing merged notes to GitHub...")
	gitrepo.AddCommitPush(notesPath, "Import offline notes")

	// Save config
	cfg, _ := config.Load()
	cfg.GitEnabled = true
	cfg.RepoURL = repo
	config.Save(cfg)

	fmt.Println("GitHub sync enabled. Offline notes merged.")
}
func dirNotEmpty(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	// read 1 entry — fast check
	_, err = f.Readdirnames(1)
	return err == nil
}

// to be implemented later
func handleConfig(args []string)    { fmt.Println("TODO: config", args) }
func handleUninstall(args []string) { fmt.Println("TODO: uninstall", args) }
