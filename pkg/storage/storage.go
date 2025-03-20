package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tejastn10/quill/pkg/hash"
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

func ReadObject(repoPath, hash string) ([]byte, error) {
	objectPath := filepath.Join(repoPath, ".quill", "objects", hash[:2], hash[2:])
	data, err := os.ReadFile(objectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}
	return data, nil
}

func WriteTree(repoPath string) (string, error) {
	var entries []string

	workDir := filepath.Join(repoPath, ".quill", "staging")
	err := filepath.Walk(workDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == workDir {
			return nil
		}

		// Compute relative path
		relPath, _ := filepath.Rel(workDir, path)

		// Read file content
		var objectHash string
		if info.IsDir() {
			// If it's a directory, recursively generate its tree object
			objectHash, err = WriteTree(repoPath)
			if err != nil {
				return err
			}
		} else {
			// Read file contents and hash it
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			objectHash = hash.ComputeSHA256(data)

			// Store as a blob
			err = CreateObject(repoPath, objectHash, data)
			if err != nil {
				return err
			}
		}

		// Format: "<hash> <filetype> <filename>"
		entry := fmt.Sprintf("%s %s %s", objectHash, "blob", relPath)
		entries = append(entries, entry)

		return nil
	})

	if err != nil {
		return "", err
	}

	// Serialize tree contents
	treeData := strings.Join(entries, "\n")
	treeHash := hash.ComputeSHA256([]byte(treeData))

	// Store tree object
	err = CreateObject(repoPath, treeHash, []byte(treeData))
	if err != nil {
		return "", err
	}

	return treeHash, nil
}
