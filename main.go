package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func createQuillRepository(path string) error {
	// Defining the Quill directory structure
	directories := []string{
		filepath.Join(path, ".quill"),
		filepath.Join(path, ".quill", "objects"),
		filepath.Join(path, ".quill", "refs"),
		filepath.Join(path, ".quill", "config"),
	}

	// Creating directories
	for _, dir := range directories {
		err := os.MkdirAll(dir, 0750)
		if err != nil {
			return fmt.Errorf("failed to create the directory %s: %w", dir, err)
		}
	}

	return nil
}

func checkQuillExists(path string) bool {
	quillPath := filepath.Join(path, ".quill")
	_, err := os.Stat(quillPath)
	return !os.IsNotExist(err)
}

func main() {
	// Get the current working Directory
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error: Unable to get the current working directory:", err)
		return
	}

	// Check if Quill is already initialized
	quillExists := checkQuillExists(workingDir)
	if quillExists {
		fmt.Println("Error: A .quill repository already exists in this directory.")
		return
	}

	// Create .quill repository structure
	err = createQuillRepository(workingDir)
	if err != nil {
		fmt.Println("Error: Failed to initialize repository:", err)
		return
	}

	fmt.Println("Initialized empty Quill repository in", filepath.Join(workingDir, ".quill"))
}
