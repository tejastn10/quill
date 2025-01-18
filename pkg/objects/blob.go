package objects

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tejastn10/quill/pkg/hash"
	"github.com/tejastn10/quill/pkg/storage"
)

// Generates a blob for the given file and stores it in the repository.
func CreateBlob(repoPath string, filePath string) (string, error) {
	// Resolve and sanitize the absolute file path
	absPath, err := filepath.Abs(filepath.Clean(filePath))
	if err != nil {
		return "", err
	}

	// Resolve and sanitize the absolute repository path
	repoAbsPath, err := filepath.Abs(filepath.Clean(repoPath))
	if err != nil {
		return "", err
	}

	// Ensure the filePath is within the repoPath
	relPath, err := filepath.Rel(repoAbsPath, absPath)
	if err != nil || relPath == ".." || filepath.IsAbs(relPath) || strings.HasPrefix(relPath, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("file path %q is outside the repository", absPath)
	}

	// Explicitly marking the resolved path as safe
	safePath := absPath

	// Reading the file contents
	// #nosec G304 - The file path is validated before being used
	data, err := os.ReadFile(safePath)
	if err != nil {
		return "", err
	}

	// Generating a hash for the file contents
	blobHash := hash.ComputeSHA256(data)

	// Storing the object in .quill/objects
	err = storage.CreateObject(repoPath, blobHash, data)
	if err != nil {
		return "", err
	}

	return blobHash, nil
}
