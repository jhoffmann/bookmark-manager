// Package service provides business logic services for the bookmark manager application.
package service

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// Folders handles folder operations like opening folders in the system file manager.
type Folders struct {
	platform string
}

// NewFolders creates a new Folders service instance.
// It caches the platform detection for better performance.
func NewFolders() *Folders {
	return &Folders{
		platform: runtime.GOOS,
	}
}

// OpenInFileManager opens the specified folder path in the system's file manager.
// Supports macOS (open), Windows (explorer), and Linux/Unix (xdg-open).
func (fs *Folders) OpenInFileManager(path string) error {
	var cmd *exec.Cmd

	switch fs.platform {
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("explorer", path)
	default: // linux and others
		cmd = exec.Command("xdg-open", path)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open folder %q: %w", path, err)
	}

	return nil
}

// WriteCwdFile writes the given path to a file for shell integration.
// This is used when the application is invoked with --cwd-file flag.
func (fs *Folders) WriteCwdFile(filePath, directoryPath string) error {
	if err := os.WriteFile(filePath, []byte(directoryPath), 0644); err != nil {
		return fmt.Errorf("failed to write to cwd file %q: %w", filePath, err)
	}
	return nil
}
