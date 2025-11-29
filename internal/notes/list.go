package notes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)

func List() error {
	root := storage.NotesDir()

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == root {
			fmt.Println("notes/")
			return nil
		}

		rel, _ := filepath.Rel(root, path)

		if info.IsDir() {
			fmt.Printf("%s%s/\n", indent(rel), info.Name())
			return nil
		}

		if strings.HasSuffix(info.Name(), ".md") {
			fmt.Printf("%s%s\n", indent(rel), info.Name())
		}

		return nil
	})
}

func indent(relPath string) string {
	parts := strings.Split(relPath, string(os.PathSeparator))
	depth := len(parts) - 1
	return strings.Repeat("    ", depth)
}

