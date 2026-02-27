package core

import (
	"io"
	"os"
	"path/filepath"
	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)

func Copy(srcRel, destRel string) error {
	notesRoot := storage.NotesDir()
	src := filepath.Join(notesRoot, srcRel)
	dest := filepath.Join(notesRoot, destRel)

	// Open source
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination
	err = os.MkdirAll(filepath.Dir(dest), 0755)
	if err != nil {
		return err
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Perform the copy
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return autoCommit("copy " + srcRel + " to " + destRel)
}
