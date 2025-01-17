package objects

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tejastn10/quill/pkg/hash"
)

func TestCreateBlob(t *testing.T) {
	// Create a temporary directory to act as the repository
	tempRepo := t.TempDir()

	// Create the .quill/objects directory
	err := os.MkdirAll(filepath.Join(tempRepo, ".quill", "objects"), os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create .quill directory: %v", err)
	}

	// Create a temporary file with test content
	tempFile := filepath.Join(tempRepo, "test.txt")
	testData := []byte("This is a test file")
	err = os.WriteFile(tempFile, testData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Call CreateBlob
	t.Run("CreateBlob", func(t *testing.T) {
		blobHash, err := CreateBlob(tempRepo, tempFile)
		if err != nil {
			t.Fatalf("CreateBlob failed: %v", err)
		}

		// Verify the hash matches the file content
		expectedHash := hash.ComputeSHA256(testData)
		if blobHash != expectedHash {
			t.Errorf("Blob hash mismatch: got %s, want %s", blobHash, expectedHash)
		}

		// Verify the blob was stored in the .quill/objects directory
		objectPath := filepath.Join(tempRepo, ".quill", "objects", blobHash[:2], blobHash[2:])
		content, err := os.ReadFile(objectPath)
		if err != nil {
			t.Fatalf("Failed to read blob object: %v", err)
		}

		if string(content) != string(testData) {
			t.Errorf("Blob content mismatch: got %q, want %q", content, testData)
		}
	})

	// Test with a non-existent file
	t.Run("NonExistentFile", func(t *testing.T) {
		nonExistentFile := filepath.Join(tempRepo, "nonexistent.txt")
		_, err := CreateBlob(tempRepo, nonExistentFile)
		if err == nil {
			t.Errorf("CreateBlob did not fail for a non-existent file")
		}
	})
}
