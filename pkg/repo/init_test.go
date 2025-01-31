package repo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateQuillRepository(t *testing.T) {
	// Creating a temporary directory for testing
	testDir := t.TempDir()

	err := CreateQuillRepository(testDir)
	if err != nil {
		t.Fatalf("Failed to create Quill repository: %v", err)
	}

	// Verifying the directory structure
	expectedDirs := []string{
		filepath.Join(testDir, ".quill"),
		filepath.Join(testDir, ".quill", "objects"),
		filepath.Join(testDir, ".quill", "refs"),
		filepath.Join(testDir, ".quill", "config"),
	}
	for _, dir := range expectedDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Expected directory %s to exist, but it doesn't", dir)
		}
	}
}

func TestCheckQuillExists(t *testing.T) {
	// Creating a temporary directory for testing
	tempDir := t.TempDir()

	// Testing if .quill does not exist
	exists := CheckQuillExists(tempDir)
	if exists {
		t.Errorf("Expected checkQuillExists to return false, but it returned true")
	}

	// Creating a .quill directory
	quillDir := filepath.Join(tempDir, ".quill")
	err := os.Mkdir(quillDir, 0750)
	if err != nil {
		t.Fatalf("Failed to create .quill directory: %v", err)
	}

	// Testing if .quill exists
	exists = CheckQuillExists(tempDir)
	if !exists {
		t.Errorf("Expected checkQuillExists to return true, but it returned false")
	}
}

func TestCreateQuillRepositoryFailsIfAlreadyExists(t *testing.T) {
	// Creating a temporary directory for testing
	tempDir := t.TempDir()

	// Creating a .quill directory manually
	quillDir := filepath.Join(tempDir, ".quill")
	err := os.MkdirAll(quillDir, 0750)
	if err != nil {
		t.Fatalf("Failed to create initial .quill directory: %v", err)
	}

	// Trying to initialize again, which should not overwrite
	err = CreateQuillRepository(tempDir)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestFindRepoRoot(t *testing.T) {
	// Creating a temporary base directory
	baseDir := t.TempDir()

	// Resolve symlinks to get the actual path (important for macOS)
	baseDir, err := filepath.EvalSymlinks(baseDir)
	if err != nil {
		t.Fatalf("Failed to resolve symlinks for baseDir: %v", err)
	}

	// Creating a nested directory structure
	nestedDir := filepath.Join(baseDir, "nested", "deep")
	err = os.MkdirAll(nestedDir, 0750)
	if err != nil {
		t.Fatalf("Failed to create nested directories: %v", err)
	}

	// Change the working directory to the nested directory
	err = os.Chdir(nestedDir)
	if err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}

	// Restore original working directory in defer
	defer func() {
		if err := os.Chdir(baseDir); err != nil {
			t.Errorf("Failed to restore original working directory: %v", err)
		}
	}()

	// Initially, there should be no repository
	_, err = FindRepoRoot()
	if err == nil {
		t.Errorf("Expected error when no .quill directory exists, but got none")
	}

	// Create a .quill directory in the base directory
	quillDir := filepath.Join(baseDir, ".quill")
	err = os.Mkdir(quillDir, 0750)
	if err != nil {
		t.Fatalf("Failed to create .quill directory: %v", err)
	}

	// Now FindRepoRoot should find the baseDir
	repoRoot, err := FindRepoRoot()
	if err != nil {
		t.Errorf("Unexpected error finding repo root: %v", err)
	}

	// Resolve symlinks before comparison
	repoRoot, err = filepath.EvalSymlinks(repoRoot)
	if err != nil {
		t.Fatalf("Failed to resolve symlinks for repoRoot: %v", err)
	}

	if repoRoot != baseDir {
		t.Errorf("Expected repo root to be %s, but got %s", baseDir, repoRoot)
	}
}
