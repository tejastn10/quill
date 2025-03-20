package repo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// CreateQuillRepository initializes a new Quill repository by creating a .quill directory structure with objects and config subdirectories
func CreateQuillRepository(path string) error {
	// Defining the Quill directory structure
	directories := []string{
		filepath.Join(path, ".quill"),
		filepath.Join(path, ".quill", "objects"),
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

// CheckQuillExists checks if a Quill repository exists at the specified path by verifying the presence of a .quill directory
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

// CleanupRepository removes the .quill directory if an error occurs
func CleanupRepository(repoPath string, err *error) {
	if *err != nil {
		quillPath := filepath.Join(repoPath, ".quill")
		fmt.Println("Rolling back: Removing partially created repository...")
		_ = os.RemoveAll(quillPath)
	}
}
