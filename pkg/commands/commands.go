package commands

import (
	"github.com/acheddir/git-contrib/pkg/stats"
)

// Stats process Git repositories and display commit statistics.
// If an email is provided, it filters commits by that email address.
// If no email is provided, it includes commits from all users.
//
// Parameters:
//   - email: The email address to filter commits by (if empty, includes all commits)
//   - directory: The directory to analyze (should be a Git repository)
//   - showCommitCount: Whether to display the number of commits on each cell
//   - showDaysOfMonth: Whether to display the days of the month on the graph calendar
//
// Returns:
//   - error: An error if any occurred during processing
func Stats(email string, directory string, showCommitCount bool, showDaysOfMonth bool) error {
	commits, err := stats.ProcessRepositories(email, directory)
	if err != nil {
		return err
	}

	stats.PrintCommitsStats(commits, showCommitCount, showDaysOfMonth)
	return nil
}
