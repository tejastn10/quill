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

	// Read HEAD
	data, err := os.ReadFile(headPath)
	if err != nil {
		return "", fmt.Errorf("failed to read HEAD: %w", err)
	}

	return strings.TrimSpace(string(data)), nil
}
