package repo

import (
	"errors"
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

// FindRepoRoot locates the root directory of the repository (with a .quill folder).
func FindRepoRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	for {
		// Check if the .quill directory exists in the current directory.
		quillPath := filepath.Join(currentDir, ".quill")
		if stat, err := os.Stat(quillPath); err == nil && stat.IsDir() {
			return currentDir, nil
		}

		// Move up one level.
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break // Reached the filesystem root without finding .quill.
		}
		currentDir = parentDir
	}

	return "", errors.New("not a quill repository (or any of the parent directories): .quill")
}
