package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tejastn10/quill/pkg/objects"
)

var addCmd = &cobra.Command{
	Use:   "add [file]",
	Short: "Add file to the repository",
	Long:  "Adds file to the .quill repository by creating its blob.",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the current working Directory
		repoPath, err := os.Getwd()
		if err != nil {
			fmt.Println("Error: Unable to get the current working directory:", err)
			return
		}

		filePath := args[0]

		// Creating blob for the given file
		hash, err := objects.CreateBlob(repoPath, filePath)
		if err != nil {
			fmt.Printf("Failed to add file: %v\n", err)
			return
		}

		fmt.Printf("File added successfully. Hash: %s\n", hash)
	},
}

func init() {
	// Registering the add command with the root command
	rootCmd.AddCommand(addCmd)
}
