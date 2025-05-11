package commands

import (
	"os"
	"path/filepath"
	"testing"
)

// TestScan tests the Scan function
func TestScan(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a .git directory to simulate a Git repository
	gitDir := filepath.Join(tempDir, ".git")
	err := os.Mkdir(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	// Call the Scan function
	// Note: This is a bit tricky to test thoroughly because it modifies the .git-contrib file
	// in the user's home directory. We'll just ensure it doesn't panic.
	Scan(tempDir)

	// In a real-world scenario, we would verify that the repository was added to the .git-contrib file,
	// but that would require mocking the fileutil package or creating a test-specific implementation.
}

// TestStats tests the Stats function
func TestStats(t *testing.T) {
	// This is difficult to test thoroughly without mocking the stats package
	// or setting up actual Git repositories with commits.
	// For now, we'll just test that it doesn't panic with an invalid email.

	// Call the Stats function with a non-existent email
	// This should not find any commits but should not error
	err := Stats("nonexistent@example.com")

	// We expect no error since the function should handle the case of no commits gracefully
	if err != nil {
		t.Errorf("Expected no error for non-existent email, got: %v", err)
	}
}

// Note: These tests are minimal and primarily ensure the functions don't panic.
// In a real-world scenario, we would use dependency injection or mocking to test
// these functions more thoroughly without relying on external dependencies.
