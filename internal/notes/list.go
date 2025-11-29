package notes

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)

type Node struct {
	Name     string
	IsFile   bool
	Children []*Node
}

func List() error {
	rootPath := storage.NotesDir()
	rootNode := &Node{Name: "notes", IsFile: false}

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == rootPath {
			return nil
		}

		rel, _ := filepath.Rel(rootPath, path)
		parts := strings.Split(rel, string(os.PathSeparator))

		insertNode(rootNode, parts, info.IsDir())
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println(rootNode.Name + "/")
	printTree(rootNode.Children, "")

	return nil
}

func insertNode(parent *Node, parts []string, isDir bool) {
	if len(parts) == 0 {
		return
	}

	name := parts[0]
	var child *Node

	for _, c := range parent.Children {
		if c.Name == name {
			child = c
			break
		}
	}

	if child == nil {
		child = &Node{
			Name:   name,
			IsFile: !isDir,
		}
		parent.Children = append(parent.Children, child)
	}

	if len(parts) > 1 {
		insertNode(child, parts[1:], isDir)
	}
}

func printTree(children []*Node, prefix string) {
	sort.Slice(children, func(i, j int) bool {
		if children[i].IsFile == children[j].IsFile {
			return children[i].Name < children[j].Name
		}
		return !children[i].IsFile // dirs first
	})

	for i, n := range children {
		isLast := i == len(children)-1
		var branch, nextPrefix string

		if isLast {
			branch = "└── "
			nextPrefix = prefix + "    "
		} else {
			branch = "├── "
			nextPrefix = prefix + "│   "
		}

		if n.IsFile {
			fmt.Println(prefix + branch + n.Name)
		} else {
			fmt.Println(prefix + branch + n.Name + "/")
			printTree(n.Children, nextPrefix)
		}
	}
}
