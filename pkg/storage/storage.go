package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

// Writing a file's contents as a blob in the .quill/objects directory.
func CreateObject(repoPath string, hash string, data []byte) error {
	// Constructing object path: .quill/objects/<first_two_hash_chars>/<rest_of_hash>
	objectDir := filepath.Join(repoPath, ".quill", "objects", hash[:2])
	objectPath := filepath.Join(objectDir, hash[2:])

	_, err := os.Stat(objectPath)
	// Object already exists
	if err == nil {
		return nil
	}

	// Creating the subdirectory if it doesn't exist
	err = os.Mkdir(objectDir, 0750) // Secure directory permissions
	if err != nil {
		return fmt.Errorf("failed to create object directory: %v", err)
	}

	// Writing the blob
	err = os.WriteFile(objectPath, data, 0600) // Secure file permissions
	if err != nil {
		return fmt.Errorf("failed to write object: %v", err)
	}

	return nil
}

func ObjectExists(repoPath string, hash string) bool {
	objectPath := filepath.Join(repoPath, ".quill", "objects", hash[:2], hash[2:])
	_, err := os.Stat(objectPath)
	return err == nil
}
