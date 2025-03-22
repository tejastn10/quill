package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tejastn10/quill/pkg/objects"
	"github.com/tejastn10/quill/pkg/repo"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit logs",
	Long:  "Display the commit history with details about changes in each commit.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Find repository root
		repoPath, err := repo.FindRepoRoot()
		if err != nil {
			return fmt.Errorf("failed to locate repository: %v", err)
		}

		// Get current HEAD commit hash
		headHash, err := repo.GetHEAD(repoPath)
		if err != nil {
			return fmt.Errorf("failed to get HEAD: %v", err)
		}

		// If no commits yet
		if headHash == "" {
			fmt.Println("No commits yet.")
			return nil
		}

		// Retrieve and display commit history
		currentHash := headHash
		for currentHash != "" {
			// Read commit
			commit, err := objects.ReadCommit(repoPath, currentHash)
			if err != nil {
				return fmt.Errorf("failed to read commit %s: %v", currentHash, err)
			}

			// Parse timestamp
			timestamp, err := time.Parse(time.RFC3339, commit.Timestamp)
			if err != nil {
				return fmt.Errorf("failed to parse timestamp: %v", err)
			}

			// Display commit header
			fmt.Printf("\033[33mcommit %s\033[0m\n", currentHash)
			fmt.Printf("Author: %s\n", commit.Author)
			fmt.Printf("Date:   %s\n\n", timestamp.Format("Mon Jan 2 15:04:05 2006 -0700"))
			fmt.Printf("    %s\n\n", commit.Message)

			// Get changes in this commit
			if commit.Parent != "" {
				changes, err := getCommitChanges(repoPath, commit.Tree, commit.Parent)
				if err != nil {
					return fmt.Errorf("failed to get commit changes: %v", err)
				}

				// Display changes
				if len(changes) > 0 {
					fmt.Println("Changes:")
					for _, change := range changes {
						fmt.Printf("    %s\n", change)
					}
					fmt.Println()
				}
			} else {
				// First commit - show all files
				files, err := objects.GetTreeFiles(repoPath, commit.Tree)
				if err != nil {
					return fmt.Errorf("failed to get files: %v", err)
				}

				fmt.Println("Files:")
				for _, file := range files {
					fmt.Printf("    %s\n", file)
				}
				fmt.Println()
			}

			// Move to parent commit
			currentHash = commit.Parent
		}

		return nil
	},
}

// getCommitChanges gets the list of changes between current and parent commits
func getCommitChanges(repoPath, currentTree, parentHash string) ([]string, error) {
	// Get parent commit
	parentCommit, err := objects.ReadCommit(repoPath, parentHash)
	if err != nil {
		return nil, fmt.Errorf("failed to read parent commit: %v", err)
	}

	// Get files from both trees
	currentFiles, err := objects.GetTreeFiles(repoPath, currentTree)
	if err != nil {
		return nil, fmt.Errorf("failed to get current files: %v", err)
	}

	parentFiles, err := objects.GetTreeFiles(repoPath, parentCommit.Tree)
	if err != nil {
		return nil, fmt.Errorf("failed to get parent files: %v", err)
	}

	// Compare files and detect changes
	currentMap := make(map[string]string)
	parentMap := make(map[string]string)

	for _, file := range currentFiles {
		parts := strings.SplitN(file, ":", 2)
		if len(parts) == 2 {
			currentMap[parts[1]] = parts[0] // path -> hash
		}
	}

	for _, file := range parentFiles {
		parts := strings.SplitN(file, ":", 2)
		if len(parts) == 2 {
			parentMap[parts[1]] = parts[0] // path -> hash
		}
	}

	// Generate changes list
	var changes []string

	// Find added and modified files
	for path, hash := range currentMap {
		parentHash, exists := parentMap[path]
		if !exists {
			changes = append(changes, fmt.Sprintf("\033[32madded:    %s\033[0m", path))
		} else if hash != parentHash {
			changes = append(changes, fmt.Sprintf("\033[33mmodified: %s\033[0m", path))
		}
	}

	// Find deleted files
	for path := range parentMap {
		if _, exists := currentMap[path]; !exists {
			changes = append(changes, fmt.Sprintf("\033[31mdeleted:  %s\033[0m", path))
		}
	}

	return changes, nil
}

func init() {
	rootCmd.AddCommand(logCmd)
}
