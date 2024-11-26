package client

import (
	"github.com/spf13/cobra"
)

var ()

func init() {
}

var installCmd = &cobra.Command{
	Use:    "install [device (e.g. /dev/sda)]",
	Short:  "Install NixOS on this machine, pre-configured for use with Nf6",
	Args:   cobra.ExactArgs(1),
	PreRun: Connect,
	Run: func(cmd *cobra.Command, args []string) {

	},
}
