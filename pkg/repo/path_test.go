package repo

import (
	"testing"
)

func TestIsPathSafe(t *testing.T) {
	tests := []struct {
		name     string
		testPath string
		want     bool
	}{
		{
			name:     "path outside repo",
			testPath: "../../../outside.txt",
			want:     false,
		},
		{
			name:     "empty path",
			testPath: "",
			want:     false,
		},
		{
			name:     "absolute path outside repo",
			testPath: "/tmp/unsafe.txt",
			want:     false,
		},
		{
			name:     "path with directory traversal",
			testPath: "../../etc/passwd",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPathSafe(tt.testPath)
			if got != tt.want {
				t.Errorf("IsPathSafe(%q) = %v, want %v", tt.testPath, got, tt.want)
			}
		})
	}
}
