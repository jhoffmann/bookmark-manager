package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	// Save original env vars
	originalDB := os.Getenv("BM_DATABASE")
	originalLog := os.Getenv("BM_LOGLEVEL")

	// Clean up after test
	defer func() {
		os.Setenv("BM_DATABASE", originalDB)
		os.Setenv("BM_LOGLEVEL", originalLog)
	}()

	tests := []struct {
		name    string
		dbEnv   string
		logEnv  string
		wantLog string
		wantErr bool
	}{
		{
			name:    "default values",
			dbEnv:   "",
			logEnv:  "",
			wantLog: "warn",
			wantErr: false,
		},
		{
			name:    "custom values",
			dbEnv:   "/custom/path.db",
			logEnv:  "info",
			wantLog: "info",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			os.Setenv("BM_DATABASE", tt.dbEnv)
			os.Setenv("BM_LOGLEVEL", tt.logEnv)

			cfg, err := Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if cfg.GetLogLevel() != tt.wantLog {
					t.Errorf("Load() log level = %v, want %v", cfg.GetLogLevel(), tt.wantLog)
				}

				if tt.dbEnv != "" && cfg.GetDatabasePath() != tt.dbEnv {
					t.Errorf("Load() database path = %v, want %v", cfg.GetDatabasePath(), tt.dbEnv)
				}

				if tt.dbEnv == "" && cfg.GetDatabasePath() == "" {
					t.Error("Load() database path should not be empty when using default")
				}
			}
		})
	}
}

func TestGetDefaultDatabasePath(t *testing.T) {
	path, err := getDefaultDatabasePath()
	if err != nil {
		t.Fatalf("getDefaultDatabasePath() error = %v", err)
	}

	if path == "" {
		t.Error("getDefaultDatabasePath() returned empty path")
	}

	// Check if path contains expected directory structure based on OS
	switch runtime.GOOS {
	case "windows":
		if !filepath.IsAbs(path) {
			t.Error("Windows path should be absolute")
		}
	case "darwin":
		if !strings.HasSuffix(filepath.Dir(path), "bookmark-manager") {
			t.Error("macOS path should end with bookmark-manager directory")
		}
	default: // Linux/Unix
		if !strings.HasSuffix(filepath.Dir(path), "bookmark-manager") {
			t.Error("Linux path should end with bookmark-manager directory")
		}
	}

	// Check if filename is correct
	if filepath.Base(path) != "bookmarks.db" {
		t.Errorf("Expected filename 'bookmarks.db', got %v", filepath.Base(path))
	}
}

func TestConfigGetters(t *testing.T) {
	cfg := &Config{
		DatabasePath: "/test/path.db",
		LogLevel:     "debug",
	}

	if cfg.GetDatabasePath() != "/test/path.db" {
		t.Errorf("GetDatabasePath() = %v, want %v", cfg.GetDatabasePath(), "/test/path.db")
	}

	if cfg.GetLogLevel() != "debug" {
		t.Errorf("GetLogLevel() = %v, want %v", cfg.GetLogLevel(), "debug")
	}
}
