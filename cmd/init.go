package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tejastn10/quill/pkg/repo"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Quill repository",
	Long:  "Create a new Quill repository by initializing .quill directory in the current directory.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the current working Directory
		workingDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("unable to get current directory: %w", err)
		}

		// Check if Quill is already exists
		quillExists := repo.CheckQuillExists(workingDir)
		if quillExists {
			return fmt.Errorf("a .quill repository already exists in this directory")
		}

		// Create .quill repository structure
		err = repo.CreateQuillRepository(workingDir)
		if err != nil {
			return fmt.Errorf("failed to initialize repository: %w", err)
		}

		// Defer cleanup in case of failure
		defer repo.CleanupRepository(workingDir, &err)

		// Ask for user details
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter your name: ")
		name, _ := reader.ReadString('\n')
		name = name[:len(name)-1] // Remove newline

		fmt.Print("Enter your email: ")
		email, _ := reader.ReadString('\n')
		email = email[:len(email)-1] // Remove newline

		// Create user config file
		err = repo.CreateUserConfig(name, email)
		if err != nil {
			return fmt.Errorf("failed to create user config file: %w", err)
		}

		// Success: Reset error before cleanup
		err = nil
		fmt.Println("Initialized empty Quill repository in", workingDir)
		return nil
	},
}

func init() {
	// Registering the init command with the root command
	rootCmd.AddCommand(initCmd)
}
