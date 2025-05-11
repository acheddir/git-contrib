package fileutil

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestGetDotfilePath tests the GetDotfilePath function
func TestGetDotfilePath(t *testing.T) {
	// Get the dotfile path
	dotfilePath := GetDotfilePath()

	// Check that the path is not empty
	if dotfilePath == "" {
		t.Error("Expected non-empty dotfile path, got empty string")
	}

	// Check that the path ends with .git-contrib
	if filepath.Base(dotfilePath) != ".git-contrib" {
		t.Errorf("Expected dotfile name to be .git-contrib, got %s", filepath.Base(dotfilePath))
	}
}

// TestParseFileLines tests the ParseFileLines function
func TestParseFileLines(t *testing.T) {
	// Create a temporary file for testing
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")

	// Test case 1: File doesn't exist
	lines := ParseFileLines(tempFile)
	if len(lines) != 0 {
		t.Errorf("Expected empty slice for non-existent file, got %v", lines)
	}

	// Test case 2: File exists with content
	content := []string{"line1", "line2", "line3"}
	err := os.WriteFile(tempFile, []byte("line1\nline2\nline3"), 0666)
	if err != nil {
		t.Fatalf("Failed to write to test file: %v", err)
	}

	lines = ParseFileLines(tempFile)
	if !reflect.DeepEqual(lines, content) {
		t.Errorf("Expected %v, got %v", content, lines)
	}

	// Test case 3: Empty file
	err = os.WriteFile(tempFile, []byte(""), 0666)
	if err != nil {
		t.Fatalf("Failed to write to test file: %v", err)
	}

	lines = ParseFileLines(tempFile)
	if len(lines) != 0 {
		t.Errorf("Expected empty slice for empty file, got %v", lines)
	}
}

// TestOpenFile tests the OpenFile function
func TestOpenFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")

	// Test opening a file that doesn't exist
	file := OpenFile(tempFile)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			t.Errorf("Failed to close file: %v", err)
		}
	}(file)

	// Check that the file was created
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Errorf("Expected file to be created, but it doesn't exist")
	}

	// Test opening a file that already exists
	file2 := OpenFile(tempFile)
	defer func(file2 *os.File) {
		err := file2.Close()
		if err != nil {
			t.Errorf("Failed to close file: %v", err)
		}
	}(file2)

	// Check that the file still exists
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Errorf("Expected file to exist, but it doesn't")
	}
}

// TestDumpStringsToFile tests the DumpStringsToFile function
func TestDumpStringsToFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")

	// Test writing to a file
	repos := []string{"repo1", "repo2", "repo3"}
	DumpStringsToFile(repos, tempFile)

	// Check that the file was created with the correct content
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	expected := "repo1\nrepo2\nrepo3"
	if string(content) != expected {
		t.Errorf("Expected content %q, got %q", expected, string(content))
	}

	// Test writing to a file in a non-existent directory
	nestedFile := filepath.Join(tempDir, "nested", "test.txt")
	DumpStringsToFile(repos, nestedFile)

	// Check that the file was created with the correct content
	content, err = os.ReadFile(nestedFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if string(content) != expected {
		t.Errorf("Expected content %q, got %q", expected, string(content))
	}
}

// TestAddElementsToFile tests the AddElementsToFile function
func TestAddElementsToFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")

	// Test adding elements to a non-existent file
	newRepos := []string{"repo1", "repo2", "repo3"}
	AddElementsToFile(tempFile, newRepos)

	// Check that the file was created with the correct content
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	expected := "repo1\nrepo2\nrepo3"
	if string(content) != expected {
		t.Errorf("Expected content %q, got %q", expected, string(content))
	}

	// Test adding more elements to the file
	moreRepos := []string{"repo3", "repo4", "repo5"}
	AddElementsToFile(tempFile, moreRepos)

	// Check that the file was updated with the correct content (no duplicates)
	content, err = os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	expected = "repo1\nrepo2\nrepo3\nrepo4\nrepo5"
	if string(content) != expected {
		t.Errorf("Expected content %q, got %q", expected, string(content))
	}
}

// TestJoinSlices tests the JoinSlices function
func TestJoinSlices(t *testing.T) {
	// Test case 1: Adding new elements to an empty slice
	newArr := []string{"a", "b", "c"}
	var existing []string
	result := JoinSlices(newArr, existing)
	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test case 2: Adding new elements to a non-empty slice
	newArr = []string{"c", "d", "e"}
	existing = []string{"a", "b", "c"}
	result = JoinSlices(newArr, existing)
	expected = []string{"a", "b", "c", "d", "e"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test case 3: Adding duplicate elements
	newArr = []string{"a", "b", "c"}
	existing = []string{"a", "b", "c"}
	result = JoinSlices(newArr, existing)
	expected = []string{"a", "b", "c"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test case 4: Adding empty slice
	newArr = []string{}
	existing = []string{"a", "b", "c"}
	result = JoinSlices(newArr, existing)
	expected = []string{"a", "b", "c"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestSliceContains tests the SliceContains function
func TestSliceContains(t *testing.T) {
	// Test case 1: Slice contains the value
	slice := []string{"a", "b", "c"}
	value := "b"
	result := SliceContains(slice, value)
	if !result {
		t.Errorf("Expected true for slice containing %q, got false", value)
	}

	// Test case 2: Slice doesn't contain the value
	value = "d"
	result = SliceContains(slice, value)
	if result {
		t.Errorf("Expected false for slice not containing %q, got true", value)
	}

	// Test case 3: Empty slice
	slice = []string{}
	value = "a"
	result = SliceContains(slice, value)
	if result {
		t.Errorf("Expected false for empty slice, got true")
	}
}
