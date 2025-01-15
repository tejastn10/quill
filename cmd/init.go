package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tejastn10/quill/pkg/repo"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Quill repository",
	Long:  "Create a new Quill repository by initializing .quill directory in the current directory.",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the current working Directory
		workingDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error: Unable to get the current working directory:", err)
			return
		}

		// Check if Quill is already initialized
		quillExists := repo.CheckQuillExists(workingDir)
		if quillExists {
			fmt.Println("Error: A .quill repository already exists in this directory.")
			return
		}

		// Create .quill repository structure
		err = repo.CreateQuillRepository(workingDir)
		if err != nil {
			fmt.Println("Error: Failed to initialize repository:", err)
			return
		}

		fmt.Println("Initialized empty Quill repository in", workingDir)
	},
}

func init() {
	// Registering the init command with the root command
	rootCmd.AddCommand(initCmd)
}
