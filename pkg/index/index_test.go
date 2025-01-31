package index

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadIndex verifies loading the index from a file.
func TestLoadIndex(t *testing.T) {
	tempDir := t.TempDir()

	// Case 1: Index file does not exist (should return an empty index)
	idx, err := LoadIndex(tempDir)
	if err != nil {
		t.Fatalf("Expected no error when index does not exist, got: %v", err)
	}
	if len(idx.Entries) != 0 {
		t.Errorf("Expected empty index, got: %v", idx.Entries)
	}

	// Case 2: Create an index and save it, then load again
	idx.Entries["file.txt"] = IndexEntry{Path: "file.txt", Hash: "abcd1234", Mode: "100644"}
	err = idx.SaveIndex(tempDir)
	if err != nil {
		t.Fatalf("Failed to save index: %v", err)
	}

	loadedIdx, err := LoadIndex(tempDir)
	if err != nil {
		t.Fatalf("Failed to load index: %v", err)
	}
	if len(loadedIdx.Entries) != 1 || loadedIdx.Entries["file.txt"].Hash != "abcd1234" {
		t.Errorf("Index data mismatch: %+v", loadedIdx.Entries)
	}
}

// TestSaveIndex verifies saving the index to a file.
func TestSaveIndex(t *testing.T) {
	tempDir := t.TempDir()

	idx := &Index{
		Entries: map[string]IndexEntry{
			"test.txt": {Path: "test.txt", Hash: "hash123", Mode: "100644"},
		},
	}

	err := idx.SaveIndex(tempDir)
	if err != nil {
		t.Fatalf("Failed to save index: %v", err)
	}

	// Ensure the file exists
	indexPath := filepath.Join(tempDir, ".quill", "index")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Fatalf("Expected index file to exist, but it doesn't")
	}

	// Load and verify content
	loadedIdx, err := LoadIndex(tempDir)
	if err != nil {
		t.Fatalf("Failed to load index: %v", err)
	}
	if loadedIdx.Entries["test.txt"].Hash != "hash123" {
		t.Errorf("Expected hash to be 'hash123', got %s", loadedIdx.Entries["test.txt"].Hash)
	}
}

// TestAddFile verifies adding files to the index.
func TestAddFile(t *testing.T) {
	tempDir := t.TempDir()
	idx := &Index{Entries: make(map[string]IndexEntry)}

	// Create a sample file
	filePath := filepath.Join(tempDir, "sample.txt")
	err := os.WriteFile(filePath, []byte("hello world"), 0644)
	if err != nil {
		t.Fatalf("Failed to create sample file: %v", err)
	}

	// Add file to index
	err = idx.AddFile(tempDir, filePath)
	if err != nil {
		t.Fatalf("Failed to add file to index: %v", err)
	}

	// Check if entry is added
	relPath, _ := filepath.Rel(tempDir, filePath)
	relPath = filepath.ToSlash(relPath)
	entry, exists := idx.Entries[relPath]
	if !exists {
		t.Fatalf("Expected file %s to be added to index", relPath)
	}
	if entry.Path != relPath {
		t.Errorf("Expected path %s, got %s", relPath, entry.Path)
	}
	if entry.Mode != "100644" {
		t.Errorf("Expected mode 100644, got %s", entry.Mode)
	}
}
