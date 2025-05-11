package fileutil

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// GetDotfilePath returns the path to the .git-contrib dotfile in the user's home directory.
// This file is used to store information about Git repositories.
//
// Returns:
//   - The full path to the .git-contrib dotfile
func GetDotfilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// Return the path to the dotfile directly in the user's home directory
	dotFile := filepath.Join(usr.HomeDir, ".git-contrib")
	return dotFile
}

// ParseFileLines reads a file and returns its contents as a slice of strings, one per line.
// If the file doesn't exist, it returns an empty slice.
//
// Parameters:
//   - filePath: The path to the file to read
//
// Returns:
//   - A slice of strings, one for each line in the file
func ParseFileLines(filePath string) []string {
	// Check if a file exists first
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Ensure the directory exists
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Failed to create directory %s: %v", dir, err)
			return []string{}
		}
		// Create an empty file
		file, err := os.Create(filePath)
		if err != nil {
			log.Printf("Failed to create file %s: %v", filePath, err)
			return []string{}
		}
		err = file.Close()
		if err != nil {
			return nil
		}
		return []string{}
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open file %s: %v", filePath, err)
		return []string{}
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Failed to close file %s: %v", filePath, err)
		}
	}(file)

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			log.Printf("Error scanning file %s: %v", filePath, err)
		}
	}

	return lines
}

// OpenFile opens a file for appending and writing, creating it if it doesn't exist.
//
// Parameters:
//   - filePath: The path to the file to open
//
// Returns:
//   - A file handle for the opened file
func OpenFile(filePath string) *os.File {
	// First, ensure the parent directory exists
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Create the directory with appropriate permissions
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Failed to create directory %s: %v", dir, err)
			panic(err)
		}
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				log.Printf("Failed to create file %s: %v", filePath, err)
				panic(err)
			}
			return file
		} else {
			log.Printf("Failed to open file %s: %v", filePath, err)
			panic(err)
		}
	}

	return file
}

// DumpStringsToFile writes a slice of strings to a file, one per line.
//
// Parameters:
//   - repos: The slice of strings to write
//   - path: The path to the file to write to
func DumpStringsToFile(repos []string, path string) {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Failed to create directory %s: %v", dir, err)
			return
		}
	}

	content := strings.Join(repos, "\n")
	err := os.WriteFile(path, []byte(content), 0666)
	if err != nil {
		log.Printf("Failed to write to file %s: %v", path, err)
		return
	}
}

// AddElementsToFile adds new elements to a file, avoiding duplicates.
//
// Parameters:
//   - filePath: The path to the file to add elements to
//   - newRepos: The new elements to add
func AddElementsToFile(filePath string, newRepos []string) {
	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Failed to create directory %s: %v", dir, err)
			return
		}
	}

	existingRepos := ParseFileLines(filePath)
	repos := JoinSlices(newRepos, existingRepos)
	DumpStringsToFile(repos, filePath)
}

// JoinSlices combines two slices, avoiding duplicates.
//
// Parameters:
//   - new: The new elements to add
//   - existing: The existing elements
//
// Return:
//   - A combined slice with no duplicates
func JoinSlices(new []string, existing []string) []string {
	for _, i := range new {
		if !SliceContains(existing, i) {
			existing = append(existing, i)
		}
	}

	return existing
}

// SliceContains checks if a slice contains a specific value.
//
// Parameters:
//   - slice: The slice to check
//   - value: The value to look for
//
// Returns:
//   - true if the slice contains the value, false otherwise
func SliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}

	return false
}
