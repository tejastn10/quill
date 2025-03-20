package objects

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/tejastn10/quill/pkg/hash"
	"github.com/tejastn10/quill/pkg/index"
	"github.com/tejastn10/quill/pkg/storage"
)

// Commit represents a commit object
type Commit struct {
	Hash      string `json:"hash"`
	Parent    string `json:"parent"`
	Timestamp string `json:"timestamp"`
	Author    string `json:"author"`
	Message   string `json:"message"`
	Tree      string `json:"tree"`
}

// CreateCommit generates a new commit from staged changes
func CreateCommit(repoPath, message, author string) (string, error) {
	// Create tree object from index
	treeHash, err := CreateTree(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to create tree: %w", err)
	}

	// Read HEAD to find last commit
	headRefPath := filepath.Join(repoPath, ".quill", "HEAD")
	headBytes, err := os.ReadFile(headRefPath)

	var parentHash string
	if err == nil {
		parentHash = string(headBytes)
	}

	// Create commit object
	commit := Commit{
		Parent:    parentHash,
		Timestamp: time.Now().Format(time.RFC3339),
		Author:    author,
		Message:   message,
		Tree:      treeHash,
	}

	// Marshal commit data
	data, err := json.Marshal(commit)
	if err != nil {
		return "", fmt.Errorf("failed to marshal commit: %w", err)
	}

	// Compute commit hash
	commit.Hash = hash.ComputeSHA256(data)

	// Update commit with hash
	data, err = json.Marshal(commit)
	if err != nil {
		return "", fmt.Errorf("failed to marshal commit with hash: %w", err)
	}

	// Store commit object
	err = storage.CreateObject(repoPath, commit.Hash, data)
	if err != nil {
		return "", fmt.Errorf("failed to store commit: %w", err)
	}

	// Update HEAD
	err = os.WriteFile(headRefPath, []byte(commit.Hash), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to update HEAD: %w", err)
	}

	// Mark all entries as unstaged
	idx, err := index.LoadIndex(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to load index: %w", err)
	}

	for path, entry := range idx.Entries {
		entry.Staged = false
		idx.Entries[path] = entry
	}

	// Save the updated index
	err = idx.SaveIndex(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to update index: %w", err)
	}

	return commit.Hash, nil
}

// ReadCommit reads a commit object from storage
func ReadCommit(repoPath, hash string) (*Commit, error) {
	// Read the commit object
	data, err := storage.ReadObject(repoPath, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to read commit: %w", err)
	}

	// Unmarshal the commit
	var commit Commit
	err = json.Unmarshal(data, &commit)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal commit: %w", err)
	}

	return &commit, nil
}

// GetTreeFiles returns a list of files in a tree
func GetTreeFiles(repoPath, treeHash string) ([]string, error) {
	// Read the tree object
	data, err := storage.ReadObject(repoPath, treeHash)
	if err != nil {
		return nil, fmt.Errorf("failed to read tree: %w", err)
	}

	var tree struct {
		Entries map[string]struct {
			Path string `json:"path"`
			Hash string `json:"hash"`
			Mode string `json:"mode"`
		} `json:"entries"`
	}

	err = json.Unmarshal(data, &tree)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tree: %w", err)
	}

	// Check if tree structure is valid
	if tree.Entries == nil {
		return nil, fmt.Errorf("invalid tree format: no entries field")
	}

	var files []string

	for path, entry := range tree.Entries {
		files = append(files, fmt.Sprintf("%s:%s", entry.Hash[:8], path))
	}

	return files, nil
}
