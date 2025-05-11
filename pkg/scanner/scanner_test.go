package scanner

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestScanFolder tests the ScanFolder function
func TestScanFolder(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Test case 1: Empty directory (no Git repositories)
	folders := ScanFolder(tempDir)
	if len(folders) != 0 {
		t.Errorf("Expected empty slice for directory with no Git repositories, got %v", folders)
	}

	// Test case 2: Directory with a Git repository
	gitDir := filepath.Join(tempDir, ".git")
	err := os.Mkdir(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	folders = ScanFolder(tempDir)
	expected := []string{tempDir}
	if !reflect.DeepEqual(folders, expected) {
		t.Errorf("Expected %v, got %v", expected, folders)
	}

	// Test case 3: Directory with nested Git repositories
	// Create a nested directory structure
	nestedDir := filepath.Join(tempDir, "nested")
	err = os.Mkdir(nestedDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create nested directory: %v", err)
	}

	nestedGitDir := filepath.Join(nestedDir, ".git")
	err = os.Mkdir(nestedGitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create nested .git directory: %v", err)
	}

	folders = ScanFolder(tempDir)
	expected = []string{tempDir, nestedDir}
	// Sort both slices to ensure consistent comparison
	if len(folders) != len(expected) {
		t.Errorf("Expected %v, got %v", expected, folders)
	} else {
		// Check that both expected repositories are in the result
		foundTemp := false
		foundNested := false
		for _, folder := range folders {
			if folder == tempDir {
				foundTemp = true
			}
			if folder == nestedDir {
				foundNested = true
			}
		}
		if !foundTemp || !foundNested {
			t.Errorf("Expected to find both %s and %s in %v", tempDir, nestedDir, folders)
		}
	}

	// Test case 4: Directory with vendor and node_modules directories (should be skipped)
	vendorDir := filepath.Join(tempDir, "vendor")
	err = os.Mkdir(vendorDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create vendor directory: %v", err)
	}

	vendorGitDir := filepath.Join(vendorDir, ".git")
	err = os.Mkdir(vendorGitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create vendor .git directory: %v", err)
	}

	nodeModulesDir := filepath.Join(tempDir, "node_modules")
	err = os.Mkdir(nodeModulesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create node_modules directory: %v", err)
	}

	nodeModulesGitDir := filepath.Join(nodeModulesDir, ".git")
	err = os.Mkdir(nodeModulesGitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create node_modules .git directory: %v", err)
	}

	folders = ScanFolder(tempDir)
	// The vendor and node_modules directories should be skipped
	if len(folders) != 2 {
		t.Errorf("Expected 2 folders (skipping vendor and node_modules), got %d: %v", len(folders), folders)
	}
}

// TestScanGitFolders tests the ScanGitFolders function
func TestScanGitFolders(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Test case 1: Empty directory (no Git repositories)
	folders := make([]string, 0)
	result := ScanGitFolders(folders, tempDir)
	if len(result) != 0 {
		t.Errorf("Expected empty slice for directory with no Git repositories, got %v", result)
	}

	// Test case 2: Directory with a Git repository
	gitDir := filepath.Join(tempDir, ".git")
	err := os.Mkdir(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	folders = make([]string, 0)
	result = ScanGitFolders(folders, tempDir)
	expected := []string{tempDir}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test case 3: Directory with nested Git repositories
	// Create a nested directory structure
	nestedDir := filepath.Join(tempDir, "nested")
	err = os.Mkdir(nestedDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create nested directory: %v", err)
	}

	nestedGitDir := filepath.Join(nestedDir, ".git")
	err = os.Mkdir(nestedGitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create nested .git directory: %v", err)
	}

	folders = make([]string, 0)
	result = ScanGitFolders(folders, tempDir)
	expected = []string{tempDir, nestedDir}
	// Sort both slices to ensure consistent comparison
	if len(result) != len(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	} else {
		// Check that both expected repositories are in the result
		foundTemp := false
		foundNested := false
		for _, folder := range result {
			if folder == tempDir {
				foundTemp = true
			}
			if folder == nestedDir {
				foundNested = true
			}
		}
		if !foundTemp || !foundNested {
			t.Errorf("Expected to find both %s and %s in %v", tempDir, nestedDir, result)
		}
	}

	// Test case 4: Directory with vendor and node_modules directories (should be skipped)
	vendorDir := filepath.Join(tempDir, "vendor")
	err = os.Mkdir(vendorDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create vendor directory: %v", err)
	}

	vendorGitDir := filepath.Join(vendorDir, ".git")
	err = os.Mkdir(vendorGitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create vendor .git directory: %v", err)
	}

	nodeModulesDir := filepath.Join(tempDir, "node_modules")
	err = os.Mkdir(nodeModulesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create node_modules directory: %v", err)
	}

	nodeModulesGitDir := filepath.Join(nodeModulesDir, ".git")
	err = os.Mkdir(nodeModulesGitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create node_modules .git directory: %v", err)
	}

	folders = make([]string, 0)
	result = ScanGitFolders(folders, tempDir)
	// The vendor and node_modules directories should be skipped
	if len(result) != 2 {
		t.Errorf("Expected 2 folders (skipping vendor and node_modules), got %d: %v", len(result), result)
	}

	// Test case 5: Existing folders are preserved and new ones are added
	existingFolders := []string{"existing/folder"}
	result = ScanGitFolders(existingFolders, tempDir)
	if len(result) != 3 { // 1 existing + 2 found
		t.Errorf("Expected 3 folders (1 existing + 2 found), got %d: %v", len(result), result)
	}

	// Check that the existing folder is still in the result
	foundExisting := false
	for _, folder := range result {
		if folder == "existing/folder" {
			foundExisting = true
			break
		}
	}
	if !foundExisting {
		t.Errorf("Expected to find existing folder in %v", result)
	}
}
