package cmd

import (
	"fmt"
	"github.com/acheddir/git-contrib/pkg/commands"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var workingDir string
var email string
var selfFlag bool
var showCommitCountFlag bool
var showDaysOfMonthFlag bool

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Display commits contribution graph",
	Long: `Process Git repositories and display commits contribution graph for all users.
This command will analyze the current working directory as a Git repository
and generate statistics about commits made by all users.
If an email is provided, it will show contributions from that email address only.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if both -c and -d flags are used together
		if showCommitCountFlag && showDaysOfMonthFlag {
			fmt.Println("Error: The -c (count) and -d (days) flags cannot be used together")
			return
		}

		// Use the specified working directory, otherwise use the current directory
		currentDir, err := filepath.Abs(workingDir)
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			return
		}

		// If the self-flag is set, get the email from git config
		if selfFlag {
			gitCmd := exec.Command("git", "config", "--global", "user.email")
			output, err := gitCmd.Output()
			if err != nil {
				fmt.Println("Error getting user email from git config:", err)
				return
			}
			email = strings.TrimSpace(string(output))
			if email == "" {
				fmt.Println("No email found in git config. Please set your email with 'git config --global user.email \"your.email@example.com\"'")
				return
			}
		}

		err = commands.Stats(email, currentDir, showCommitCountFlag, showDaysOfMonthFlag)
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)

	// Add the working directory flag to the stats command
	statsCmd.Flags().StringVarP(&workingDir, "path", "p", ".", "The directory to analyze (default is the current working directory)")

	// Add the email flag to the stats command (no default value)
	statsCmd.Flags().StringVarP(&email, "email", "e", "", "The email address to filter commits by (if empty, shows all users)")

	// Add the self-flag to use the current user's email from git config
	statsCmd.Flags().BoolVarP(&selfFlag, "self", "s", false, "Use the current user's email from git config")

	// Add flags to show the commit count on cells and days of the month
	statsCmd.Flags().BoolVarP(&showCommitCountFlag, "count", "c", false, "Display the number of commits on each cell")
	statsCmd.Flags().BoolVarP(&showDaysOfMonthFlag, "days", "d", false, "Display the days of the month on the graph calendar")

	// Make stats the default command when no subcommand is specified
	cobra.OnInitialize(func() {
		// If no subcommand is specified, run the stats command
		if len(os.Args) == 1 {
			statsCmd.Run(statsCmd, []string{})
			os.Exit(0)
		}
	})
}
