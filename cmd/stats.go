package cmd

import (
	"fmt"
	"git-contrib/pkg/commands"
	"github.com/spf13/cobra"
	"os"
)

var email string

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Display commit statistics for a given email",
	Long: `Process Git repositories and display commit statistics for a given email.
This command will analyze all repositories in the .git-contrib dotfile
and generate statistics about commits made by the specified email address.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := commands.Stats(email)
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)

	// Add the email flag to the stats command
	statsCmd.Flags().StringVarP(&email, "email", "e", "acheddir@redsen.ch", "The email address to filter commits by")

	// Make stats the default command when no subcommand is specified
	cobra.OnInitialize(func() {
		// If no subcommand is specified, run the stats command
		if len(os.Args) == 1 {
			statsCmd.Run(statsCmd, []string{})
			os.Exit(0)
		}
	})
}
