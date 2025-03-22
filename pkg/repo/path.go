package repo

import (
	"path/filepath"
	"strings"
)

func IsPathSafe(path string) bool {
	basePath, err := FindRepoRoot()
	if err != nil {
		return false
	}

	// Get absolute paths
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	absBasePath, err := filepath.Abs(basePath)
	if err != nil {
		return false
	}

	// Check if the path is within the base path
	return strings.HasPrefix(absPath, absBasePath)
}
