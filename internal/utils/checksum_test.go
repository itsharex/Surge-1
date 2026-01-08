package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"testing"
)

func TestVerifyChecksum_MD5(t *testing.T) {
	// Create a temp file with known content
	content := []byte("Hello, Surge!")
	tmpFile, err := os.CreateTemp("", "checksum_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Calculate expected MD5
	hasher := md5.New()
	hasher.Write(content)
	expectedMD5 := hex.EncodeToString(hasher.Sum(nil))

	// Test matching MD5
	err = VerifyChecksum(tmpFile, expectedMD5, "", "", "", false)
	if err != nil {
		t.Errorf("VerifyChecksum failed for correct MD5: %v", err)
	}

	// Test mismatching MD5
	err = VerifyChecksum(tmpFile, "wrongchecksum", "", "", "", false)
	if err == nil {
		t.Error("VerifyChecksum should fail for incorrect MD5")
	}
}

func TestVerifyChecksum_SHA256(t *testing.T) {
	// Create a temp file with known content
	content := []byte("Hello, Surge!")
	tmpFile, err := os.CreateTemp("", "checksum_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Calculate expected SHA256
	hasher := sha256.New()
	hasher.Write(content)
	expectedSHA256 := hex.EncodeToString(hasher.Sum(nil))

	// Test matching SHA256
	err = VerifyChecksum(tmpFile, "", expectedSHA256, "", "", false)
	if err != nil {
		t.Errorf("VerifyChecksum failed for correct SHA256: %v", err)
	}

	// Test mismatching SHA256
	err = VerifyChecksum(tmpFile, "", "wrongchecksum", "", "", false)
	if err == nil {
		t.Error("VerifyChecksum should fail for incorrect SHA256")
	}
}

func TestVerifyChecksum_NoChecksum(t *testing.T) {
	// Create a temp file
	tmpFile, err := os.CreateTemp("", "checksum_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	tmpFile.Write([]byte("test content"))

	// When no checksum provided, should pass (skip verification)
	err = VerifyChecksum(tmpFile, "", "", "", "", false)
	if err != nil {
		t.Errorf("VerifyChecksum should pass when no checksum provided: %v", err)
	}
}

func TestVerifyChecksum_ServerMD5(t *testing.T) {
	content := []byte("Server checksum test")
	tmpFile, err := os.CreateTemp("", "checksum_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Calculate expected MD5
	hasher := md5.New()
	hasher.Write(content)
	expectedMD5 := hex.EncodeToString(hasher.Sum(nil))

	// Test with server-provided MD5 (3rd param is serverMD5)
	err = VerifyChecksum(tmpFile, "", "", expectedMD5, "", false)
	if err != nil {
		t.Errorf("VerifyChecksum failed for correct server MD5: %v", err)
	}
}
