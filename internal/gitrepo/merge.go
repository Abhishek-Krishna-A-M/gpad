package gitrepo

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// MergeOfflineIntoRepo copies offline notes into the repo with conflict prompts.
func MergeOfflineIntoRepo(repoDir, offlineDir string) error {
	return filepath.Walk(offlineDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || path == offlineDir {
			return err
		}
		rel, _ := filepath.Rel(offlineDir, path)
		target := filepath.Join(repoDir, rel)

		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		if _, err := os.Stat(target); err == nil {
			switch askConflict(rel) {
			case 1:
				return nil
			case 2:
				return copyFile(path, target)
			case 3:
				ext := filepath.Ext(rel)
				base := rel[:len(rel)-len(ext)]
				_ = os.Rename(target, filepath.Join(repoDir, base+"_remote"+ext))
				return copyFile(path, filepath.Join(repoDir, base+"_local"+ext))
			}
			return nil
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	_ = os.MkdirAll(filepath.Dir(dst), 0755)
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func askConflict(name string) int {
	fmt.Printf("Conflict: %s\n  1 keep remote  2 keep local  3 keep both\n> ", name)
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _ := reader.ReadString('\n')
		switch strings.TrimSpace(line) {
		case "1":
			return 1
		case "2":
			return 2
		case "3":
			return 3
		default:
			fmt.Print("Choose 1, 2, or 3: ")
		}
	}
}
