package client

import (
	"github.com/spf13/cobra"
)

var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate a new TLS keypair",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
