package client

import (
	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func init() {
	accountCmd.AddCommand(accountGetCmd)
	accountCmd.AddCommand(accountSetEmailCmd)
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

var accountSetEmailCmd = &cobra.Command{
	Use:    "set-email [email]",
	Short:  "Set your account email",
	Args:   cobra.MaximumNArgs(1),
	PreRun: Connect,
	Run: func(cmd *cobra.Command, args []string) {
		email := ""
		if len(args) > 0 {
			email = args[0]
		}
		if err := lib.PromptOrValidate(&email, &promptui.Prompt{
			Label:    "Email",
			Validate: lib.ValidateEmail,
		}); err != nil {
			lib.Crash(err)
		}
		ctx, cancel := lib.Context()
		defer cancel()
		_, err := client.UpdateAccount(ctx, &nf6.UpdateAccount_Request{Email: &email})
		if err != nil {
			lib.Crash(err)
		}
	},
}
