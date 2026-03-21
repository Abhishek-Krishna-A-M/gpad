package notes

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/config"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)

const (
	bold   = "\033[1m"
	cyan   = "\033[96m"
	yellow = "\033[93m"
	reset  = "\033[0m"
	dim    = "\033[2m"
)

type node struct {
	name     string
	isFile   bool
	children []*node
}

// List prints the vault as a tree with pinned markers.
func List() error {
	root := &node{name: "notes", isFile: false}
	rootPath := storage.NotesDir()

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		if path == rootPath {
			return nil
		}
		rel, _ := filepath.Rel(rootPath, path)
		parts := strings.Split(rel, string(os.PathSeparator))
		insert(root, parts, info.IsDir())
		return nil
	})
	if err != nil {
		return err
	}

	fmt.Printf("%s%snotes/%s\n", bold, cyan, reset)
	printTree(root.children, "", rootPath)
	return nil
}

func insert(parent *node, parts []string, isDir bool) {
	if len(parts) == 0 {
		return
	}
	name := parts[0]
	var child *node
	for _, c := range parent.children {
		if c.name == name {
			child = c
			break
		}
	}
	if child == nil {
		child = &node{name: name, isFile: !isDir}
		parent.children = append(parent.children, child)
	}
	if len(parts) > 1 {
		insert(child, parts[1:], isDir)
	}
}

func printTree(children []*node, prefix, rootPath string) {
	sort.Slice(children, func(i, j int) bool {
		if children[i].isFile == children[j].isFile {
			return children[i].name < children[j].name
		}
		return !children[i].isFile
	})
	for i, n := range children {
		isLast := i == len(children)-1
		branch := "├── "
		nextPrefix := prefix + "│   "
		if isLast {
			branch = "└── "
			nextPrefix = prefix + "    "
		}
		if n.isFile {
			// compute relative path from notesRoot for pin check
			rel := filepath.Join(strings.TrimPrefix(prefix, ""), n.name)
			pin := ""
			if config.IsPinned(rel) {
				pin = yellow + " ★" + reset
			}
			fmt.Printf("%s%s%s%s%s\n", prefix, branch, n.name, pin, reset)
		} else {
			fmt.Printf("%s%s%s%s/%s\n", prefix, branch, bold, n.name, reset)
			printTree(n.children, nextPrefix, rootPath)
		}
	}
}
