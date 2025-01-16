package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateObjectAndObjectExists(t *testing.T) {
	// Create a temporary directory to act as the repository
	tempDir := t.TempDir()

	// Create the .quill directory
	err := os.MkdirAll(filepath.Join(tempDir, ".quill", "objects"), os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create .quill directory: %v", err)
	}

	// Initialize test data
	repoPath := tempDir
	hash := "abcd1234efgh5678ijkl9101mnopqrst"
	data := []byte("This is a test blob")

	// Test CreateObject
	t.Run("CreateObject", func(t *testing.T) {
		err := CreateObject(repoPath, hash, data)
		if err != nil {
			t.Fatalf("CreateObject failed: %v", err)
		}

		// Verify the object file was created with the correct content
		objectPath := filepath.Join(repoPath, ".quill", "objects", hash[:2], hash[2:])
		content, err := os.ReadFile(objectPath)
		if err != nil {
			t.Fatalf("Failed to read object file: %v", err)
		}

		if string(content) != string(data) {
			t.Errorf("Object content mismatch: got %q, want %q", content, data)
		}
	})

	// Test ObjectExists
	t.Run("ObjectExists", func(t *testing.T) {
		exists := ObjectExists(repoPath, hash)
		if !exists {
			t.Errorf("ObjectExists returned false for an existing object")
		}

		// Test with a non-existent object
		nonExistentHash := "0000111122223333444455556666777788889999"
		exists = ObjectExists(repoPath, nonExistentHash)
		if exists {
			t.Errorf("ObjectExists returned true for a non-existent object")
		}
	})
}
