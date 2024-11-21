package client

import (
	"github.com/computerdane/nf6/lib"
	"github.com/spf13/cobra"
)

var (
	gentlsGenCa   bool
	gentlsGenCert bool
	gentlsDir     string
	gentlsName    string
	gentlsCaName  string
)

func init() {
	gentlsCmd.PersistentFlags().BoolVar(&gentlsGenCa, "ca", false, "generate a ca cert & keypair")
	gentlsCmd.PersistentFlags().BoolVar(&gentlsGenCert, "cert", false, "generate a cert")
	gentlsCmd.PersistentFlags().StringVarP(&gentlsDir, "dir", "d", "", "directory to put new keypair")
	gentlsCmd.PersistentFlags().StringVarP(&gentlsName, "name", "n", "", "file name for keypair")
	gentlsCmd.PersistentFlags().StringVarP(&gentlsCaName, "ca-name", "", "", "file name for ca cert")
}

var gentlsCmd = &cobra.Command{
	Use:   "gentls",
	Short: "Generate a new TLS keypair",
	Run: func(cmd *cobra.Command, args []string) {
		if gentlsDir == "" {
			gentlsDir = tlsDir
		}
		if gentlsCaName == "" {
			gentlsCaName = tlsCaName
		}
		if gentlsName == "" {
			gentlsName = tlsName
		}
		if gentlsGenCa {
			if err := lib.GenCaFiles(gentlsDir, gentlsCaName); err != nil {
				lib.Crash(err)
			}
		} else if gentlsGenCert {
			if err := lib.GenCertFiles(gentlsDir, gentlsCaName, gentlsName); err != nil {
				lib.Crash(err)
			}
		} else {
			if err := lib.GenKeyFiles(gentlsDir, gentlsName); err != nil {
				lib.Crash(err)
			}
		}
	},
}
