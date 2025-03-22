package repo

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tejastn10/quill/pkg/constants"
)

// CreateUserConfig creates the user's config when quill is initialized
func CreateUserConfig(name string, email string) error {
	if !isValidEmail(email) {
		return fmt.Errorf("invalid email format in config")
	}

	// Get the current working Directory
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to get the current working directory: %w", err)
	}

	// Define and sanitize the config directory path
	userConfigDir := filepath.Join(workingDir, ".quill", "config")
	userConfigDir = filepath.Clean(userConfigDir)

	// Ensure the config directory is inside workingDir
	if !strings.HasPrefix(userConfigDir, workingDir) {
		return fmt.Errorf("invalid user config directory path: %s", userConfigDir)
	}

	// Create user config directory
	err = os.MkdirAll(userConfigDir, constants.DirectoryPerms)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Define and sanitize the user config file path
	userConfigFile := filepath.Join(userConfigDir, "user")
	userConfigFile = filepath.Clean(userConfigFile)

	// Ensure the config file is inside the expected directory
	if !strings.HasPrefix(userConfigFile, userConfigDir) {
		return fmt.Errorf("invalid user config file path: %s", userConfigFile)
	}

	// Create and write to the file with explicit permissions
	file, err := os.OpenFile(userConfigFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, constants.ConfigFilePerms)
	if err != nil {
		return fmt.Errorf("failed to create user config file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("name=%s\nemail=%s\n", name, email))
	if err != nil {
		return fmt.Errorf("failed to write user config file: %w", err)
	}

	return nil
}

// ReadUserConfig reads the user's name and email from the config file
func ReadUserConfig(repoPath string) (string, string, error) {
	configPath := filepath.Join(repoPath, ".quill", "config", "user")

	cleanPath := filepath.Clean(configPath)

	if !IsPathSafe(cleanPath) {
		return "", "", fmt.Errorf("invalid file path: potential directory traversal attempt")
	}

	file, err := os.Open(cleanPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to open user config: %w", err)
	}
	defer file.Close()

	var name, email string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "name":
			name = value
		case "email":
			email = value
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", fmt.Errorf("error reading config file: %w", err)
	}

	if name == "" || email == "" {
		return "", "", fmt.Errorf("user name or email not found in config")
	}

	return name, email, nil
}

// isValidEmail checks if the given email follows a valid format
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
