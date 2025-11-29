package cli

import (
	"fmt"
	"os"
	"strings"
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
func dirExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func extractUser(url string) string {
	// https://github.com/user/repo.git → user
	parts := strings.Split(url, "/")
	if len(parts) < 5 {
		return ""
	}
	return parts[3]
}

func extractRepo(url string) string {
	// https://github.com/user/repo.git → repo.git → repo
	parts := strings.Split(url, "/")
	if len(parts) < 5 {
		return ""
	}
	repo := parts[len(parts)-1]
	return strings.TrimSuffix(repo, ".git")
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
	var repoURL string

	if len(args) >= 2 && args[0] == "--github" {
		repoURL = args[1]
	}

	notesPath := storage.NotesDir()

	// OFFLINE MODE
	if repoURL == "" {
		fmt.Println("Initializing gpad in offline mode...")
		fmt.Println("Notes stored at:", notesPath)
		return
	}

	// GIT MODE
	fmt.Println("Initializing gpad with GitHub repo:", repoURL)

	if !gitrepo.Exists() {
		fmt.Println("Error: git is not installed.")
		return
	}

	isSSH := strings.HasPrefix(repoURL, "git@github.com:")
	isHTTPS := strings.HasPrefix(repoURL, "https://github.com/")

	// Warn if HTTPS
	if isHTTPS {
		fmt.Println("⚠️  Using HTTPS for GitHub sync.")
		fmt.Println("You may be prompted for PAT (GitHub token).")
		fmt.Println("")
		fmt.Println("Recommended SSH mode:")
		fmt.Printf("  git@github.com:%s/%s.git\n", extractUser(repoURL), extractRepo(repoURL))
		fmt.Println("")
		fmt.Println("Or enable credential caching:")
		fmt.Println("  git config --global credential.helper store")
		fmt.Println("")
	}

	// If repo already exists & is git repo → update remote (don't clone again)
	if dirExists(filepath.Join(notesPath, ".git")) {
		if isSSH {
			fmt.Println("Switching existing repo to SSH mode:", repoURL)
			if err := gitrepo.SetRemote(notesPath, repoURL); err != nil {
				fmt.Println("Failed to set SSH remote:", err)
				return
			}
		}
		cfg, _ := config.Load()
		cfg.GitEnabled = true
		cfg.RepoURL = repoURL
		config.Save(cfg)
		fmt.Println("GitHub sync is now active.")
		return
	}

	// FRESH CLONE FLOW
	tmpDir := filepath.Join(storage.GpadDir(), "tmp_clone")

	// Clean temp dir
	os.RemoveAll(tmpDir)
	os.MkdirAll(tmpDir, 0755)

	// Clone to tmp
	if err := gitrepo.Clone(repoURL, tmpDir); err != nil {
		fmt.Println("Git clone failed:", err)
		return
	}

	// Merge offline notes
	if dirNotEmpty(notesPath) {
		fmt.Println("Merging offline notes into GitHub repo...")
		if err := gitrepo.MergeOfflineIntoRepo(tmpDir, notesPath); err != nil {
			fmt.Println("Merge failed:", err)
			return
		}
	}

	// Replace notes folder
	os.RemoveAll(notesPath)
	os.Rename(tmpDir, notesPath)

	// Auto commit merged notes
	fmt.Println("Syncing merged notes to GitHub...")
	gitrepo.AddCommitPush(notesPath, "Import offline notes")

	// Save config
	cfg, _ := config.Load()
	cfg.GitEnabled = true
	cfg.RepoURL = repoURL
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

func handleConfig(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage:")
		fmt.Println("  gpad config editor <command>")
		fmt.Println("  gpad config autopush on/off")
		return
	}

	if args[0] == "editor" && len(args) >= 2 {
		cmd := strings.Join(args[1:], " ")

		cfg, _ := config.Load()
		cfg.Editor = cmd
		config.Save(cfg)

		fmt.Println("Editor set to:", cmd)
		return
	}



	if args[0] == "autopush" && len(args) == 2 {
		cfg, _ := config.Load()
		if args[1] == "on" {
			cfg.AutoPush = true
		} else {
			cfg.AutoPush = false
		}
		config.Save(cfg)
		fmt.Println("autopush set to:", cfg.AutoPush)
		return
	}

	fmt.Println("Unknown config command")
}

// to be implemented later
func handleUninstall(args []string) { fmt.Println("TODO: uninstall", args) }
