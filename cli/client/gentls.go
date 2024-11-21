package client

import (
	"github.com/computerdane/nf6/lib"
	"github.com/spf13/cobra"
)

var (
	gentlsGenCa bool
	gentlsDir   string
	gentlsName  string
)

func init() {
	gentlsCmd.PersistentFlags().BoolVar(&gentlsGenCa, "ca", false, "generate a ca cert & keypair")
	gentlsCmd.PersistentFlags().StringVarP(&gentlsDir, "dir", "d", "", "directory to put new keypair")
	gentlsCmd.PersistentFlags().StringVarP(&gentlsName, "name", "n", "", "file name for keypair")
}

var gentlsCmd = &cobra.Command{
	Use:   "gentls",
	Short: "Generate a new TLS keypair",
	Run: func(cmd *cobra.Command, args []string) {
		if gentlsDir == "" {
			gentlsDir = tlsDir
		}
		if gentlsGenCa {
			if gentlsName == "" {
				gentlsName = tlsCaName
			}
			if err := lib.GenCaFiles(gentlsDir, gentlsName); err != nil {
				lib.Crash(err)
			}
		} else {
			if gentlsName == "" {
				gentlsName = tlsName
			}
			if err := lib.GenKeyFiles(gentlsDir, gentlsName); err != nil {
				lib.Crash(err)
			}
		}
	},
}
