package utils

import (
	"testing"
)

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple filename", "file.zip", "file.zip"},
		{"filename with spaces", "  file.zip  ", "file.zip"},
		{"filename with backslash", "path\\file.zip", "file.zip"},
		{"filename with forward slash", "path/file.zip", "file.zip"},
		{"filename with colon", "file:name.zip", "file_name.zip"},
		{"filename with asterisk", "file*name.zip", "file_name.zip"},
		{"filename with question mark", "file?name.zip", "file_name.zip"},
		{"filename with quotes", "file\"name.zip", "file_name.zip"},
		{"filename with angle brackets", "file<name>.zip", "file_name_.zip"},
		{"filename with pipe", "file|name.zip", "file_name.zip"},
		{"dot only", ".", "."},
		// filepath.Base("/") returns "/" on Unix but "\" on Windows, then sanitized to "_"
		{"slash only", "/", "_"},
		// filepath.Base extracts "d.zip" from "a:b*c?d.zip" on Windows (treats a: as drive)
		{"multiple bad chars", "b*c?d.zip", "b_c_d.zip"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeFilename(tt.input)
			if got != tt.expected {
				t.Errorf("sanitizeFilename(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestConvertBytesToHumanReadable(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
		{1099511627776, "1.0 TB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := ConvertBytesToHumanReadable(tt.bytes)
			if got != tt.expected {
				t.Errorf("ConvertBytesToHumanReadable(%d) = %q, want %q", tt.bytes, got, tt.expected)
			}
		})
	}
}
