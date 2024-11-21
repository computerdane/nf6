package client

import (
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register [email]",
	Short: "Register using your email",
	Args:  cobra.ExactArgs(1),

	PreRun: ConnectPublic,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
