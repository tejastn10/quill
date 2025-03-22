package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStorageFunctions(t *testing.T) {
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
}
