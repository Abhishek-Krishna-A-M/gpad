package storage

import (
	"os"
	"path/filepath"
)

func HomeDir() string {
	home, _ := os.UserHomeDir()
	return home
}

func GpadDir() string {
	return filepath.Join(HomeDir(), ".gpad")
}

func NotesDir() string {
	return filepath.Join(GpadDir(), "notes")
}

func EnsureDirs() error {
	paths := []string{
		GpadDir(),
		NotesDir(),
	}

	for _, p := range paths {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			if err := os.MkdirAll(p, 0755); err != nil {
				return err
			}
		}
	}
	return nil
}
func AbsPath(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}
	return filepath.Join(NotesDir(), rel)
}
