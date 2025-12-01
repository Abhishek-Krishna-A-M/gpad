package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)
func handleRename(args []string) {
	if len(args) != 2 {
		fmt.Println("Usage: gpad mv <old> <new>")
		return
	}

	oldRel := args[0]
	newRel := args[1]

	notesRoot := storage.NotesDir()

	oldPath := filepath.Join(notesRoot, oldRel)
	newPath := filepath.Join(notesRoot, newRel)

	oldPath = filepath.Clean(oldPath)
	newPath = filepath.Clean(newPath)

	if !strings.HasPrefix(oldPath, notesRoot) || !strings.HasPrefix(newPath, notesRoot) {
		fmt.Println("Error: rename only allowed inside", notesRoot)
		return
	}

	if oldPath == notesRoot {
		fmt.Println("Error: cannot rename the notes root directory")
		return
	}

	info, err := os.Stat(oldPath)
	if err != nil {
		fmt.Println("Error: file or directory not found:", oldRel)
		return
	}

	if _, err := os.Stat(newPath); err == nil {
		fmt.Println("Error: destination already exists:", newRel)
		return
	}

	os.MkdirAll(filepath.Dir(newPath), 0755)

	err = os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Println("Rename failed:", err)
		return
	}

	if !info.IsDir() {
		updateHeading(newPath)
	}

	fmt.Printf("Renamed '%s' → '%s'\n", oldRel, newRel)
	commitAndPush("rename " + oldRel + " to " + newRel)
}
func updateHeading(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")

	base := filepath.Base(path)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	name = strings.ToLower(strings.ReplaceAll(name, "_", " "))
	name = strings.ReplaceAll(name, "-", " ")

	for i, line := range lines {
		trim := strings.TrimSpace(line)

		if strings.HasPrefix(trim, "# ") {
			lines[i] = "# " + name
			goto WRITE
		}

		if trim == "" {
			continue
		}

		break
	}

	return

WRITE:
	newContent := strings.Join(lines, "\n")
	_ = os.WriteFile(path, []byte(newContent), 0644)
}

