package repo

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tejastn10/quill/pkg/constants"
)

func TestCreateUserConfig(t *testing.T) {
	// Create a temporary directory using t.TempDir()
	tempDir := t.TempDir()

	// Change the working directory to the temporary directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("Failed to restore original directory: %v", err)
		}
	}()

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}

	// Call the function with test data
	name := "Test User"
	email := "test@example.com"
	err = CreateUserConfig(name, email)
	if err != nil {
		t.Fatalf("CreateUserConfig returned an error: %v", err)
	}

	// Verify the user config file was created with the correct content
	userConfigFile := filepath.Join(tempDir, ".quill", "config", "user")
	content, err := os.ReadFile(userConfigFile)
	if err != nil {
		t.Fatalf("Failed to read user config file: %v", err)
	}

	expectedContent := "name=Test User\nemail=test@example.com\n"
	if string(content) != expectedContent {
		t.Errorf("User config file content mismatch. Expected:\n%s\nGot:\n%s", expectedContent, string(content))
	}
}

func TestReadUserConfig_InvalidEmail(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Create necessary directories
	configDir := filepath.Join(tempDir, ".quill", "config")
	err := os.MkdirAll(configDir, constants.DirectoryPerms)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	// Define the config file path
	userConfigFile := filepath.Join(configDir, "user")

	// Write an invalid email format
	content := "name=Test User\nemail=invalid-email\n"
	err = os.WriteFile(userConfigFile, []byte(content), constants.ConfigFilePerms)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Call ReadUserConfig and verify it returns an error
	_, _, err = ReadUserConfig(tempDir)
	if err == nil {
		t.Fatal("Expected ReadUserConfig to return an error for invalid email, but it didn't")
	}
}

func TestReadUserConfig_MissingConfigFile(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Call ReadUserConfig and verify it returns an error
	_, _, err := ReadUserConfig(tempDir)
	if err == nil {
		t.Fatal("Expected ReadUserConfig to return an error for missing config file, but it didn't")
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name+tag@example.co.uk", true},
		{"user@sub.domain.com", true},
		{"invalid-email", false},
		{"@missingusername.com", false},
		{"username@.nodomain", false},
		{"username@domain..com", false},
		{"plainaddress", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := isValidEmail(tt.email)
			if result != tt.expected {
				t.Errorf("isValidEmail(%s) = %v; want %v", tt.email, result, tt.expected)
			}
		})
	}
}
