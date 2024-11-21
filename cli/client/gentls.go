package client

import (
	"github.com/computerdane/nf6/lib"
	"github.com/spf13/cobra"
)

var (
	gentlsDir string
)

func init() {
	gentlsCmd.PersistentFlags().StringVarP(&gentlsDir, "dir", "d", "", "directory to put new keypair")
}

var gentlsCmd = &cobra.Command{
	Use:   "gentls",
	Short: "Generate a new TLS keypair",
	Run: func(cmd *cobra.Command, args []string) {
		if gentlsDir == "" {
			gentlsDir = tlsDir
		}
		if err := lib.GenKeyFiles(gentlsDir, "client"); err != nil {
			lib.Crash(err)
		}
	},
}
