package notes

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/editor"
)

func Create(relPath string) error {
	base := storage.NotesDir()
	full := filepath.Join(base, relPath)

	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return err
	}

	if _, err := os.Stat(full); err == nil {
		return fmt.Errorf("note already exists: %s", relPath)
	}

	f, err := os.Create(full)
	if err != nil {
		return err
	}
	defer f.Close()

	title := filepath.Base(relPath)
	noExt := title[:len(title)-len(filepath.Ext(title))]

	timestamp := time.Now().Format("2006-01-02 15:04")

	header := fmt.Sprintf("# %s\n%s\n\n", noExt, timestamp)

	if _, err := f.WriteString(header); err != nil {
		return err
	}

	return editor.Open(full)
}

