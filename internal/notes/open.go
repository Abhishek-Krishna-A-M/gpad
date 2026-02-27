package notes

import (
	"os"
	"path/filepath"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/editor"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/core"
)

func Open(relPath string) error {
	// This ensures you are never editing an outdated version of a note
	_ = core.Sync() 

	base := storage.NotesDir()
	full := filepath.Join(base, relPath)

	// 2. Open or Create
	if _, err := os.Stat(full); os.IsNotExist(err) {
		if err := Create(relPath); err != nil {
			return err
		}
	} else {
		if err := editor.Open(full); err != nil {
			return err
		}
	}

	// 3. Background Sync: Only pushes if AutoPush is enabled in config
	core.AutoSave("Update " + relPath)
	
	return nil
}
