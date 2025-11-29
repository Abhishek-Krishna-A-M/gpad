package gitrepo

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"bufio"
)

func MergeOfflineIntoRepo(repoDir, offlineDir string) error {
	return filepath.Walk(offlineDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == offlineDir {
			return nil
		}

		rel, _ := filepath.Rel(offlineDir, path)
		target := filepath.Join(repoDir, rel)

		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		if fileExists(target) {
			choice := askConflict(rel)

			switch choice {
			case 1: // keep remote (do nothing)
				return nil

			case 2: // keep offline (overwrite)
				return copyFile(path, target)

			case 3: // keep both → rename
				ext := filepath.Ext(rel)
				name := rel[:len(rel)-len(ext)]

				localPath := filepath.Join(repoDir, name+"_local"+ext)
				remotePath := filepath.Join(repoDir, name+"_remote"+ext)

				os.Rename(target, remotePath)

				return copyFile(path, localPath)
			}
			return nil
		}

		return copyFile(path, target)
	})
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func copyFile(src, dst string) error {
	// ensure folder
	os.MkdirAll(filepath.Dir(dst), 0755)

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
	fmt.Println("Conflict:", name)
	fmt.Println("  1 = keep GitHub version")
	fmt.Println("  2 = keep offline version")
	fmt.Println("  3 = keep both (local/remote)")
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	for {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		switch line {
		case "1", "2", "3":
			return int(line[0] - '0')
		default:
			fmt.Print("Choose 1, 2, or 3: ")
		}
	}
}

