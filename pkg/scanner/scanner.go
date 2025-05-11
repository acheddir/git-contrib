package scanner

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// ScanGitFolders scans the specified folder for Git repositories and adds them to the existing list of folders.
// It recursively searches through the directory structure looking for .git directories.
//
// Parameters:
//   - folders: The existing list of Git repository folders
//   - folder: The path to the folder to scan for Git repositories
//
// Returns:
//   - An updated list of Git repository folders including any new ones found
func ScanGitFolders(folders []string, folder string) []string {
	// Ensure the folder path uses the correct separator for the OS
	folder = filepath.Clean(folder)

	folderOpen, folderOpenErr := os.Open(folder)
	if folderOpenErr != nil {
		log.Fatal(folderOpenErr)
	}

	files, filesReadErr := folderOpen.Readdir(-1)
	if filesReadErr != nil {
		log.Fatal(filesReadErr)
	}

	folderCloseErr := folderOpen.Close()
	if folderCloseErr != nil {
		log.Fatal(folderCloseErr)
	}

	var path string

	for _, file := range files {
		if file.IsDir() {
			path = filepath.Join(folder, file.Name())
			if file.Name() == ".git" {
				path = filepath.Dir(path) // Remove the .git part
				fmt.Println(path)
				folders = append(folders, path)
				continue
			}
			if file.Name() == "vendor" || file.Name() == "node_modules" {
				continue
			}
			folders = ScanGitFolders(folders, path)
		}
	}

	return folders
}

// ScanFolder initializes an empty slice and calls ScanGitFolders to scan the specified folder.
//
// Parameters:
//   - folder: The path to the folder to scan for Git repositories
//
// Returns:
//   - A list of Git repository folders found
func ScanFolder(folder string) []string {
	return ScanGitFolders(make([]string, 0), folder)
}
