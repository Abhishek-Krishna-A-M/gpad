package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/notes"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
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

// to be implemented later
func handleView(path string)        { fmt.Println("TODO: view", path) }
func handleList()                   { fmt.Println("TODO: list") }
func handleInit(args []string)      { fmt.Println("TODO: init", args) }
func handleConfig(args []string)    { fmt.Println("TODO: config", args) }
func handleUninstall(args []string) { fmt.Println("TODO: uninstall", args) }
