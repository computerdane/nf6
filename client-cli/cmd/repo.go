package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	repoCmd.AddCommand(repoCreateCmd)
	rootCmd.AddCommand(repoCmd)
}

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage your repos",
}

var repoCreateCmd = &cobra.Command{
	Use:    "create [name]",
	Short:  "Create a repo",
	Args:   cobra.ExactArgs(1),
	PreRun: requireSecureClient,
	Run: func(cmd *cobra.Command, args []string) {
		// name := args[0]
	},
}
