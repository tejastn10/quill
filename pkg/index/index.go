package index

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tejastn10/quill/pkg/hash"
)

// IndexEntry represents a single entry in the index file.
type IndexEntry struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
	Mode string `json:"mode"`
}

// Index represents the staging area.
type Index struct {
	Entries map[string]IndexEntry `json:"entries"`
}

// LoadIndex loads the index from the .quill/index file.
func LoadIndex(repoPath string) (*Index, error) {
	indexPath := filepath.Join(repoPath, ".quill", "index")
	indexPath = filepath.Clean(indexPath) // Clean the path to remove potential traversal issues

	// Ensure the index file is within the repo
	if !strings.HasPrefix(indexPath, filepath.Clean(repoPath)) {
		return nil, fmt.Errorf("index path %q is outside the repository", indexPath)
	}

	file, err := os.Open(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If the index file doesn't exist, return a new empty index.
			return &Index{Entries: make(map[string]IndexEntry)}, nil
		}
		return nil, fmt.Errorf("failed to open index: %w", err)
	}
	defer file.Close()

	var idx Index
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&idx); err != nil {
		return nil, fmt.Errorf("failed to decode index: %w", err)
	}
	return &idx, nil
}

// SaveIndex saves the index to the .quill/index file.
func (idx *Index) SaveIndex(repoPath string) error {
	indexPath := filepath.Join(repoPath, ".quill", "index")
	indexPath = filepath.Clean(indexPath)

	// Ensure the .quill directory exists.
	err := os.MkdirAll(filepath.Dir(indexPath), 0750)
	if err != nil {
		return fmt.Errorf("failed to create directory for index: %w", err)
	}

	// Ensure the index file is within the repo
	if !strings.HasPrefix(indexPath, filepath.Clean(repoPath)) {
		return fmt.Errorf("index path %q is outside the repository", indexPath)
	}

	// Write the index to the file safely
	file, err := os.OpenFile(indexPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create index file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print the JSON for debugging.
	if err := encoder.Encode(idx); err != nil {
		return fmt.Errorf("failed to encode index: %w", err)
	}
	return nil
}

// AddFile adds a single file to the index.
func (idx *Index) AddFile(repoPath, filePath string) error {
	// Compute the absolute path and ensure it's within the repo
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path for %q: %w", filePath, err)
	}
	absPath = filepath.Clean(absPath)

	// Ensure the file is inside the repository
	if !strings.HasPrefix(absPath, filepath.Clean(repoPath)+string(os.PathSeparator)) {
		return fmt.Errorf("file path %q is outside the repository", filePath)
	}

	// Compute the hash of the file's contents.
	data, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	fileHash := hash.ComputeSHA256(data)

	// Compute the relative path for the file.
	relPath, err := filepath.Rel(repoPath, absPath)
	if err != nil || strings.HasPrefix(relPath, "..") {
		return fmt.Errorf("file path %q is outside the repository", filePath)
	}

	// Normalize the path and add to the index.
	relPath = filepath.ToSlash(relPath)
	idx.Entries[relPath] = IndexEntry{
		Path: relPath,
		Hash: fileHash,
		Mode: "100644", // Default mode for regular files.
	}
	return nil
}
