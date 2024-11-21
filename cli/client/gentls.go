package client

import (
	"github.com/computerdane/nf6/lib"
	"github.com/spf13/cobra"
)

var gentlsCmd = &cobra.Command{
	Use:   "gentls",
	Short: "Generate a new TLS private key and certificate",
	Run: func(cmd *cobra.Command, args []string) {
		if err := lib.GenKeyFiles(tlsDir, "client"); err != nil {
			lib.Crash(err)
		}
	},
}
