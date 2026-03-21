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

func DailyDir() string {
	return filepath.Join(NotesDir(), "daily")
}

func TemplatesDir() string {
	return filepath.Join(GpadDir(), "templates")
}

func IndexPath() string {
	return filepath.Join(GpadDir(), "index.json")
}

func EnsureDirs() error {
	paths := []string{
		GpadDir(),
		NotesDir(),
		DailyDir(),
		TemplatesDir(),
	}
	for _, p := range paths {
		if err := os.MkdirAll(p, 0755); err != nil {
			return err
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
