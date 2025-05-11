package cmd

import (
	"git-contrib/pkg/commands"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [folder]",
	Short: "Scan a folder for Git repositories",
	Long: `Scan a folder for Git repositories and add them to the .git-contrib dotfile.
This command will recursively search the specified folder for Git repositories
and add them to the list of repositories to be analyzed.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		folder := args[0]
		commands.Scan(folder)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
