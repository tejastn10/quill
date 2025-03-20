package objects

import (
	"encoding/json"
	"fmt"

	"github.com/tejastn10/quill/pkg/hash"
	"github.com/tejastn10/quill/pkg/index"
	"github.com/tejastn10/quill/pkg/storage"
)

// TreeEntry represents an entry in a tree object
type TreeEntry struct {
	Mode string `json:"mode"`
	Type string `json:"type"`
	Hash string `json:"hash"`
	Path string `json:"path"`
}

// Tree represents a tree object which contains references to blobs and other trees
type Tree struct {
	Entries []TreeEntry `json:"entries"`
}

// CreateTree creates a tree object from the current index
func CreateTree(repoPath string) (string, error) {
	// Load index
	idx, err := index.LoadIndex(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to load index: %w", err)
	}

	// Create tree object
	treeData, err := json.Marshal(idx)
	if err != nil {
		return "", fmt.Errorf("failed to marshal tree: %w", err)
	}

	// Compute tree hash
	treeHash := hash.ComputeSHA256(treeData)

	// Store tree object
	err = storage.CreateObject(repoPath, treeHash, treeData)
	if err != nil {
		return "", fmt.Errorf("failed to store tree object: %w", err)
	}

	return treeHash, nil
}
