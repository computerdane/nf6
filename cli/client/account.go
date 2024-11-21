package client

import (
	"github.com/computerdane/nf6/lib"
	"github.com/spf13/cobra"
)

func init() {
	accountCmd.AddCommand(accountGetCmd)
}

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage your account",
}

var accountGetCmd = &cobra.Command{
	Use:    "get",
	Short:  "Get your account info",
	PreRun: Connect,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := lib.Context()
		defer cancel()
		reply, err := client.GetAccount(ctx, nil)
		if err != nil {
			lib.Crash(err)
		}
		lib.Output(reply)
	},
}
