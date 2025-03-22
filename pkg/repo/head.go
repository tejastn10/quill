package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetHEAD returns the current HEAD commit hash
func GetHEAD(repoPath string) (string, error) {
	headPath := filepath.Join(repoPath, ".quill", "HEAD")

	// Check if HEAD exists
	if _, err := os.Stat(headPath); os.IsNotExist(err) {
		return "", nil // No commits yet
	}

	cleanPath := filepath.Clean(headPath)

	if !IsPathSafe(cleanPath) {
		return "", fmt.Errorf("invalid file path: potential directory traversal attempt")
	}
	// Read HEAD
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to read HEAD: %w", err)
	}

	return strings.TrimSpace(string(data)), nil
}
