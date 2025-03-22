package index

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tejastn10/quill/pkg/constants"
	"github.com/tejastn10/quill/pkg/hash"
	"github.com/tejastn10/quill/pkg/repo"
	"github.com/tejastn10/quill/pkg/storage"
)

// IndexEntry represents a single entry in the index file.
type IndexEntry struct {
	Path   string `json:"path"`
	Hash   string `json:"hash"`
	Mode   string `json:"mode"`
	Staged bool   `json:"staged,omitempty"`
}

// Index represents the staging area.
type Index struct {
	Entries        map[string]IndexEntry `json:"entries"`
	LastCommitTree string                `json:"lastCommitTree,omitempty"`
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
	err := os.MkdirAll(filepath.Dir(indexPath), constants.DirectoryPerms)
	if err != nil {
		return fmt.Errorf("failed to create directory for index: %w", err)
	}

	// Ensure the index file is within the repo
	if !strings.HasPrefix(indexPath, filepath.Clean(repoPath)) {
		return fmt.Errorf("index path %q is outside the repository", indexPath)
	}

	// Write the index to the file safely
	file, err := os.OpenFile(indexPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, constants.ConfigFilePerms)
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

// AddFile adds a file to the index with its current hash
func (idx *Index) AddFile(repoPath, filePath string) error {
	// Get file info
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to stat file %q: %w", filePath, err)
	}

	// Check if it's a regular file
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%q is not a regular file", filePath)
	}

	cleanPath := filepath.Clean(filePath)

	if !repo.IsPathSafe(cleanPath) {
		return fmt.Errorf("invalid file path: potential directory traversal attempt")
	}

	// Read file content
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", filePath, err)
	}

	// Compute hash
	fileHash := hash.ComputeSHA256(data)

	// Get relative path for storage
	relPath, err := filepath.Rel(repoPath, filePath)
	if err != nil {
		return fmt.Errorf("failed to get relative path for %q: %w", filePath, err)
	}

	// Check if file has changed since last commit
	currentEntry, exists := idx.Entries[relPath]
	if exists && currentEntry.Hash == fileHash && !currentEntry.Staged {
		// File hasn't changed, no need to add it again
		fmt.Printf("File %q unchanged, not adding to staging area\n", relPath)
		return nil
	}

	// Store the object
	err = storage.CreateObject(repoPath, fileHash, data)
	if err != nil {
		return fmt.Errorf("failed to store object for %q: %w", filePath, err)
	}

	// Add to index
	mode := fmt.Sprintf("%o", info.Mode().Perm())
	idx.Entries[relPath] = IndexEntry{
		Path:   relPath,
		Hash:   fileHash,
		Mode:   mode,
		Staged: true, // Mark as staged
	}

	fmt.Printf("Added %q to staging area\n", relPath)
	return nil
}

// CreateCleanIndex creates a clean index after commit, marking all files as committed
func CreateCleanIndex(repoPath, treeHash string) error {
	// Load the current index
	idx, err := LoadIndex(repoPath)
	if err != nil {
		return fmt.Errorf("failed to load index: %w", err)
	}

	// Store the tree hash for comparison in future operations
	// This helps track which files have changed since the last commit
	idx.LastCommitTree = treeHash

	// Keep entries but mark them as committed (useful for status command)
	for path, entry := range idx.Entries {
		entry.Staged = false // Add this field to your IndexEntry struct
		idx.Entries[path] = entry
	}

	// Save the updated index
	return idx.SaveIndex(repoPath)
}
