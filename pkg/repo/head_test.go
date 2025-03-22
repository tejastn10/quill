package repo

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tejastn10/quill/pkg/constants"
)

func TestGetHEAD(t *testing.T) {
	// Create a temporary repository
	tempDir := t.TempDir()
	headPath := filepath.Join(tempDir, ".quill", "HEAD")

	tests := []struct {
		name        string
		setup       func()
		want        string
		expectError bool
	}{
		{
			name: "HEAD does not exist",
			setup: func() {
				// Ensure .quill directory exists, but no HEAD file
				os.MkdirAll(filepath.Join(tempDir, ".quill"), os.ModePerm)
			},
			want:        "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got, err := GetHEAD(tempDir)

			if (err != nil) != tt.expectError {
				t.Fatalf("GetHEAD() error = %v, expectError %v", err, tt.expectError)
			}

			if got != tt.want {
				t.Errorf("GetHEAD() = %q, want %q", got, tt.want)
			}

			// Cleanup permission changes
			os.Chmod(headPath, constants.ConfigFilePerms)
		})
	}
}
