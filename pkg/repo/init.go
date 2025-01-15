package repo

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateQuillRepository(path string) error {
	// Defining the Quill directory structure
	directories := []string{
		filepath.Join(path, ".quill"),
		filepath.Join(path, ".quill", "objects"),
		filepath.Join(path, ".quill", "refs"),
		filepath.Join(path, ".quill", "config"),
	}

	// Creating directories
	for _, dir := range directories {
		err := os.MkdirAll(dir, 0750)
		if err != nil {
			return fmt.Errorf("failed to create the directory %s: %w", dir, err)
		}
	}

	return nil
}

func CheckQuillExists(path string) bool {
	quillPath := filepath.Join(path, ".quill")
	_, err := os.Stat(quillPath)
	return !os.IsNotExist(err)
}
