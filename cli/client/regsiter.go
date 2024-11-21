package client

import (
	"github.com/spf13/cobra"
)

func init() {
	initCommand(RegisterCmd)
}

var RegisterCmd = &cobra.Command{
	Use:   "register [email]",
	Short: "Register using your email",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
