package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tejastn10/quill/pkg/index"
	"github.com/tejastn10/quill/pkg/repo"
)

var addCmd = &cobra.Command{
	Use:   "add [files...]",
	Short: "Add file contents to the staging area",
	Long:  "Add file contents to the staging area to be included in the next commit.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Locate the repository root.
		repoPath, err := repo.FindRepoRoot()
		if err != nil {
			return fmt.Errorf("failed to locate repository: %v", err)
		}

		// Load the index.
		idx, err := index.LoadIndex(repoPath)
		if err != nil {
			return fmt.Errorf("failed to load index: %v", err)
		}

		// Process each file or directory.
		for _, arg := range args {
			// Resolve the absolute path.
			absPath, err := filepath.Abs(arg)
			if err != nil {
				return fmt.Errorf("failed to resolve path for %q: %v", arg, err)
			}

			// Add file or directory to the index.
			info, err := os.Stat(absPath)
			if err != nil {
				return fmt.Errorf("failed to stat %q: %v", absPath, err)
			}

			if info.IsDir() {
				// Recursively add files in the directory.
				err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}

					if !info.IsDir() {
						err = idx.AddFile(repoPath, path)
						if err != nil {
							return fmt.Errorf("failed to add %q: %v", path, err)
						}
					}

					return nil
				})
				if err != nil {
					return fmt.Errorf("failed to add directory %q: %v", absPath, err)
				}
			} else {
				// Add a single file.
				err = idx.AddFile(repoPath, absPath)
				if err != nil {
					return fmt.Errorf("failed to add %q: %v", absPath, err)
				}
			}
		}

		// Save the updated index.
		err = idx.SaveIndex(repoPath)
		if err != nil {
			return fmt.Errorf("failed to save index: %v", err)
		}

		fmt.Println("Files have been added to the staging area.")
		return nil
	},
}

func init() {
	// Registering the add command with the root command
	rootCmd.AddCommand(addCmd)
}
