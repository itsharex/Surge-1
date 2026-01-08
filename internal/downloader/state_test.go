package downloader

import (
	"testing"

	"surge/internal/config"
)

func TestURLHash(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantLen int
	}{
		{"simple URL", "https://example.com/file.zip", 16},
		{"URL with path", "https://example.com/path/to/file.zip", 16},
		{"URL with query", "https://example.com/file.zip?token=abc", 16},
		{"different domain", "https://other.org/download", 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := URLHash(tt.url)
			if len(hash) != tt.wantLen {
				t.Errorf("URLHash(%s) length = %d, want %d", tt.url, len(hash), tt.wantLen)
			}
		})
	}
}

func TestURLHashUniqueness(t *testing.T) {
	url1 := "https://example.com/file1.zip"
	url2 := "https://example.com/file2.zip"

	hash1 := URLHash(url1)
	hash2 := URLHash(url2)

	if hash1 == hash2 {
		t.Errorf("Different URLs produced same hash: %s", hash1)
	}
}

func TestURLHashConsistency(t *testing.T) {
	url := "https://example.com/consistent.zip"

	hash1 := URLHash(url)
	hash2 := URLHash(url)

	if hash1 != hash2 {
		t.Errorf("Same URL produced different hashes: %s vs %s", hash1, hash2)
	}
}

func TestSaveLoadState(t *testing.T) {
	// Ensure directories exist
	if err := config.EnsureDirs(); err != nil {
		t.Fatalf("Failed to create directories: %v", err)
	}

	testURL := "https://test.example.com/save-load-test.zip"
	originalState := &DownloadState{
		URL:        testURL,
		DestPath:   "C:\\Downloads\\testfile.zip",
		TotalSize:  1000000,
		Downloaded: 500000,
		Tasks: []Task{
			{Offset: 500000, Length: 250000},
			{Offset: 750000, Length: 250000},
		},
		Filename: "save-load-test.zip",
	}

	// Save state
	if err := SaveState(testURL, originalState); err != nil {
		t.Fatalf("SaveState failed: %v", err)
	}

	// Load state
	loadedState, err := LoadState(testURL)
	if err != nil {
		t.Fatalf("LoadState failed: %v", err)
	}

	// Verify fields
	if loadedState.URL != originalState.URL {
		t.Errorf("URL = %s, want %s", loadedState.URL, originalState.URL)
	}
	if loadedState.Downloaded != originalState.Downloaded {
		t.Errorf("Downloaded = %d, want %d", loadedState.Downloaded, originalState.Downloaded)
	}
	if loadedState.TotalSize != originalState.TotalSize {
		t.Errorf("TotalSize = %d, want %d", loadedState.TotalSize, originalState.TotalSize)
	}
	if len(loadedState.Tasks) != len(originalState.Tasks) {
		t.Errorf("Tasks count = %d, want %d", len(loadedState.Tasks), len(originalState.Tasks))
	}
	if loadedState.Filename != originalState.Filename {
		t.Errorf("Filename = %s, want %s", loadedState.Filename, originalState.Filename)
	}

	// Verify URLHash was set
	if loadedState.URLHash == "" {
		t.Error("URLHash was not set")
	}

	// Cleanup
	if err := DeleteState(testURL); err != nil {
		t.Errorf("Cleanup failed: %v", err)
	}
}

func TestDeleteState(t *testing.T) {
	if err := config.EnsureDirs(); err != nil {
		t.Fatalf("Failed to create directories: %v", err)
	}

	testURL := "https://test.example.com/delete-test.zip"
	state := &DownloadState{
		URL:      testURL,
		Filename: "delete-test.zip",
	}

	// Save state
	if err := SaveState(testURL, state); err != nil {
		t.Fatalf("SaveState failed: %v", err)
	}

	// Verify it was saved
	if _, err := LoadState(testURL); err != nil {
		t.Fatalf("State was not saved properly: %v", err)
	}

	// Delete state
	if err := DeleteState(testURL); err != nil {
		t.Fatalf("DeleteState failed: %v", err)
	}

	// Verify it was deleted
	_, err := LoadState(testURL)
	if err == nil {
		t.Error("LoadState should fail after DeleteState")
	}
}

func TestStateOverwrite(t *testing.T) {
	// This tests the user's scenario: pause at 30%, resume to 80%, pause again
	// The state should reflect 80%, not 30%
	if err := config.EnsureDirs(); err != nil {
		t.Fatalf("Failed to create directories: %v", err)
	}

	testURL := "https://test.example.com/overwrite-test.zip"

	// First pause at 30%
	state1 := &DownloadState{
		URL:        testURL,
		TotalSize:  1000000,
		Downloaded: 300000, // 30%
		Tasks:      []Task{{Offset: 300000, Length: 700000}},
		Filename:   "overwrite-test.zip",
	}
	if err := SaveState(testURL, state1); err != nil {
		t.Fatalf("First SaveState failed: %v", err)
	}

	// Second pause at 80% (simulating resume + more downloading)
	state2 := &DownloadState{
		URL:        testURL,
		TotalSize:  1000000,
		Downloaded: 800000, // 80%
		Tasks:      []Task{{Offset: 800000, Length: 200000}},
		Filename:   "overwrite-test.zip",
	}
	if err := SaveState(testURL, state2); err != nil {
		t.Fatalf("Second SaveState failed: %v", err)
	}

	// Load and verify it's 80%, not 30%
	loaded, err := LoadState(testURL)
	if err != nil {
		t.Fatalf("LoadState failed: %v", err)
	}

	if loaded.Downloaded != 800000 {
		t.Errorf("Downloaded = %d, want 800000 (state should be overwritten)", loaded.Downloaded)
	}
	if len(loaded.Tasks) != 1 || loaded.Tasks[0].Offset != 800000 {
		t.Errorf("Tasks not properly overwritten, got offset %d", loaded.Tasks[0].Offset)
	}

	// Cleanup
	DeleteState(testURL)
}
