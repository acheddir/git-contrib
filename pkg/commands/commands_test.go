package commands

import (
	"os"
	"path/filepath"
	"testing"
)

// TestStats tests the Stats function
func TestStats(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a .git directory to simulate a Git repository
	gitDir := filepath.Join(tempDir, ".git")
	err := os.Mkdir(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	// Call the Stats function with a non-existent email
	// This should not find any commits but should not error
	err = Stats("nonexistent@example.com", tempDir, false, false)

	// We expect an error since the directory is not a valid Git repository
	if err == nil {
		t.Errorf("Expected an error for invalid repository, got nil")
	}
}

// Note: These tests are minimal and primarily ensure the functions don't panic.
// In a real-world scenario, we would use dependency injection or mocking to test
// these functions more thoroughly without relying on external dependencies.
