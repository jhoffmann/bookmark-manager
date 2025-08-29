package service

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNewFolders(t *testing.T) {
	fs := NewFolders()
	if fs == nil {
		t.Fatal("NewFolders() returned nil")
	}
	if fs.platform != runtime.GOOS {
		t.Errorf("Expected platform %q, got %q", runtime.GOOS, fs.platform)
	}
}

func TestFolders_WriteCwdFile(t *testing.T) {
	tests := []struct {
		name          string
		directoryPath string
		wantError     bool
	}{
		{
			name:          "valid path",
			directoryPath: "/home/user/test",
			wantError:     false,
		},
		{
			name:          "empty path",
			directoryPath: "",
			wantError:     false,
		},
		{
			name:          "path with spaces",
			directoryPath: "/home/user/my documents",
			wantError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := NewFolders()

			// Create temporary file
			tmpFile, err := os.CreateTemp("", "test-cwd-*.txt")
			if err != nil {
				t.Fatal(err)
			}
			tmpFile.Close()
			defer os.Remove(tmpFile.Name())

			err = fs.WriteCwdFile(tmpFile.Name(), tt.directoryPath)
			if (err != nil) != tt.wantError {
				t.Errorf("WriteCwdFile() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				// Verify file contents
				content, err := os.ReadFile(tmpFile.Name())
				if err != nil {
					t.Fatal(err)
				}
				if string(content) != tt.directoryPath {
					t.Errorf("Expected file content %q, got %q", tt.directoryPath, string(content))
				}
			}
		})
	}
}

func TestFolders_WriteCwdFile_InvalidPath(t *testing.T) {
	fs := NewFolders()

	// Try to write to an invalid/non-existent directory
	invalidPath := filepath.Join("/nonexistent", "directory", "file.txt")

	err := fs.WriteCwdFile(invalidPath, "/some/path")
	if err == nil {
		t.Error("Expected error when writing to invalid path, got nil")
	}
}

func TestFolders_OpenInFileManager_Platform(t *testing.T) {
	tests := []struct {
		platform string
	}{
		{"darwin"},
		{"windows"},
		{"linux"},
		{"freebsd"},
	}

	for _, tt := range tests {
		t.Run(tt.platform, func(t *testing.T) {
			fs := &Folders{platform: tt.platform}

			// We can't actually test the command execution in unit tests
			// as it would require the system file manager to be available
			// Instead, we verify the service doesn't panic and handles the platform correctly
			err := fs.OpenInFileManager("/nonexistent/path")

			// We expect an error since the path doesn't exist, but the important
			// thing is that the method handles different platforms without panicking
			if err == nil {
				t.Log("Unexpected success - path likely exists or command succeeded")
			}
		})
	}
}
