package commands

import (
	"fmt"
	"git-contrib/pkg/fileutil"
	"git-contrib/pkg/scanner"
	"git-contrib/pkg/stats"
)

// Scan scans a folder for Git repositories and adds them to the .git-contrib dotfile.
//
// Parameters:
//   - folder: The path to the folder to scan for Git repositories
func Scan(folder string) {
	fmt.Printf("Found folders:\n\n")
	repos := scanner.ScanFolder(folder)
	filePath := fileutil.GetDotfilePath()
	fileutil.AddElementsToFile(filePath, repos)
	fmt.Printf("\n\nSuccessfully added\n\n")
}

// Stats process Git repositories and display commit statistics for a given email.
//
// Parameters:
//   - email: The email address to filter commits by
//
// Returns:
//   - error: An error if any occurred during processing
func Stats(email string) error {
	commits, err := stats.ProcessRepositories(email)
	if err != nil {
		return err
	}

	stats.PrintCommitsStats(commits)
	return nil
}
