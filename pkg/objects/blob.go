package objects

import (
	"os"

	"github.com/tejastn10/quill/pkg/hash"
	"github.com/tejastn10/quill/pkg/storage"
)

// Generates a blob for the given file and stores it in the repository.
func CreateBlob(repoPath string, filePath string) (string, error) {
	// Reading the file contents
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Generating a hash for the file contents
	blobHash := hash.ComputeSHA1(data)

	// Storing the object in .quill/objects
	err = storage.CreateObject(repoPath, blobHash, data)
	if err != nil {
		return "", err
	}

	return blobHash, nil
}
