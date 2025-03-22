package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tejastn10/quill/pkg/index"
	"github.com/tejastn10/quill/pkg/objects"
	"github.com/tejastn10/quill/pkg/repo"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Record changes to the repository",
	Long:  "Create a new commit containing the current contents of the index and the given log message describing the changes.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get commit message
		message, err := cmd.Flags().GetString("message")
		if err != nil {
			return fmt.Errorf("failed to get message flag: %v", err)
		}

		if message == "" {
			return fmt.Errorf("commit message cannot be empty")
		}

		// Find repository root
		repoPath, err := repo.FindRepoRoot()
		if err != nil {
			return fmt.Errorf("failed to locate repository: %v", err)
		}

		// Load the index
		idx, err := index.LoadIndex(repoPath)
		if err != nil {
			return fmt.Errorf("failed to load index: %v", err)
		}

		// Check if there are any staged changes
		hasChanges := false
		for _, entry := range idx.Entries {
			if entry.Staged {
				hasChanges = true
				break
			}
		}

		if !hasChanges {
			return fmt.Errorf("no changes staged for commit")
		}

		// Get user info from config
		name, email, err := repo.ReadUserConfig(repoPath)
		if err != nil {
			return fmt.Errorf("failed to read user config: %v", err)
		}

		// Create author string
		author := fmt.Sprintf("%s <%s>", name, email)

		// Create commit object
		commitHash, err := objects.CreateCommit(repoPath, message, author)
		if err != nil {
			return fmt.Errorf("failed to create commit: %v", err)
		}

		fmt.Printf("Created commit %s: %s\n", commitHash[:8], message)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringP("message", "m", "", "Commit message")
	err := commitCmd.MarkFlagRequired("message")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marking flag as required: %v\n", err)
		os.Exit(1)
	}
}
