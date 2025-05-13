package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information
var (
	// Version is the current version of git-contrib
	Version = "1.0.2"

	// BuildDate is the date when the binary was built
	BuildDate = "undefined"

	// CommitHash is the git commit hash when the binary was built
	CommitHash = "undefined"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version of git-contrib",
	Long:  `Display the version, build date, and commit hash of git-contrib.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("git-contrib version %s\n", Version)

		if BuildDate != "undefined" {
			fmt.Printf("Built on %s\n", BuildDate)
		}

		if CommitHash != "undefined" {
			fmt.Printf("Commit %s\n", CommitHash)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
