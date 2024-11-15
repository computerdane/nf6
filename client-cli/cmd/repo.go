package cmd

import (
	"context"

	"github.com/computerdane/nf6/nf6"
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
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		reply, err := clientSecure.CreateRepo(ctx, &nf6.CreateRepoRequest{Name: args[0]})
		if err != nil {
			crash(err)
		}
		if !reply.GetSuccess() {
			crash()
		}
	},
}
