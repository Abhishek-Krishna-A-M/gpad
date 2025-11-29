package cli

import (
	"fmt"
	"os"
)

func Run() {
	args := os.Args[1:]

	if len(args) == 0 {
		printHelp()
		return
	}

	switch args[0] {

	case "create":
		if len(args) < 2 {
			fmt.Println("Usage: gpad create <path>")
			return
		}
		handleCreate(args[1])

	case "view":
		if len(args) < 2 {
			fmt.Println("Usage: gpad view <file.md>")
			return
		}
		handleView(args[1])

	case "edit":
		if len(args) < 2 {
			fmt.Println("Usage: gpad edit <file.md>")
			return
		}
		handleEdit(args[1])

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

Commands:
  gpad init [--github URL]   Initialize gpad global notes
  gpad create <path>         Create a new note
  gpad edit <path>           Edit an existing note
  gpad view <path>           View a markdown file
  gpad list                  List all notes
  gpad config ...            Configuration commands
  gpad uninstall [...]       Uninstall gpad
`)
}

// placeholders for future logic
func handleCreate(path string)      { fmt.Println("TODO: create", path) }
func handleView(path string)        { fmt.Println("TODO: view", path) }
func handleEdit(path string)        { fmt.Println("TODO: edit", path) }
func handleList()                   { fmt.Println("TODO: list") }
func handleInit(args []string)      { fmt.Println("TODO: init", args) }
func handleConfig(args []string)    { fmt.Println("TODO: config", args) }
func handleUninstall(args []string) { fmt.Println("TODO: uninstall", args) }
