package config

import (
	"os"
	"strings"
	"testing"
)

func TestGetSurgeDir(t *testing.T) {
	dir := GetSurgeDir()
	if dir == "" {
		t.Error("GetSurgeDir returned empty string")
	}
	// Should contain "surge" in path
	if !strings.Contains(strings.ToLower(dir), "surge") {
		t.Errorf("Expected path to contain 'surge', got: %s", dir)
	}
}

func TestGetStateDir(t *testing.T) {
	dir := GetStateDir()
	if dir == "" {
		t.Error("GetStateDir returned empty string")
	}
	if !strings.HasSuffix(dir, "state") {
		t.Errorf("Expected path to end with 'state', got: %s", dir)
	}
	// State dir should be under surge dir
	surgeDir := GetSurgeDir()
	if !strings.HasPrefix(dir, surgeDir) {
		t.Errorf("StateDir should be under SurgeDir. StateDir: %s, SurgeDir: %s", dir, surgeDir)
	}
}

func TestGetLogsDir(t *testing.T) {
	dir := GetLogsDir()
	if dir == "" {
		t.Error("GetLogsDir returned empty string")
	}
	if !strings.HasSuffix(dir, "logs") {
		t.Errorf("Expected path to end with 'logs', got: %s", dir)
	}
	// Logs dir should be under surge dir
	surgeDir := GetSurgeDir()
	if !strings.HasPrefix(dir, surgeDir) {
		t.Errorf("LogsDir should be under SurgeDir. LogsDir: %s, SurgeDir: %s", dir, surgeDir)
	}
}

func TestEnsureDirs(t *testing.T) {
	err := EnsureDirs()
	if err != nil {
		t.Fatalf("EnsureDirs failed: %v", err)
	}

	// Verify all directories exist
	dirs := []string{GetSurgeDir(), GetStateDir(), GetLogsDir()}
	for _, dir := range dirs {
		info, err := os.Stat(dir)
		if os.IsNotExist(err) {
			t.Errorf("Directory not created: %s", dir)
		} else if err != nil {
			t.Errorf("Error checking directory %s: %v", dir, err)
		} else if !info.IsDir() {
			t.Errorf("Path exists but is not a directory: %s", dir)
		}
	}
}
