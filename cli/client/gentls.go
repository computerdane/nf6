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
		gentlsCaPrivKeyPath := gentlsDir + "/" + gentlsCaName + ".key"
		gentlsCaPubKeyPath := gentlsDir + "/" + gentlsCaName + ".pub"
		gentlsCaCertPath := gentlsDir + "/" + gentlsCaName + ".crt"
		gentlsPrivKeyPath := gentlsDir + "/" + gentlsName + ".key"
		gentlsPubKeyPath := gentlsDir + "/" + gentlsName + ".pub"
		gentlsCertPath := gentlsDir + "/" + gentlsName + ".crt"
		if gentlsGenCa {
			caPubKey, caPrivKey, err := lib.TlsGenKey()
			if err != nil {
				lib.Crash(err)
			}
			caCert, err := lib.TlsGenCert(lib.TlsCaTemplate, caPubKey, caPrivKey)
			if err != nil {
				lib.Crash(err)
			}
			caPubKeyPem, err := lib.TlsEncodePubKey(caPubKey)
			if err != nil {
				lib.Crash(err)
			}
			caPrivKeyPem, err := lib.TlsEncodePrivKey(caPrivKey)
			if err != nil {
				lib.Crash(err)
			}
			if err := lib.TlsWriteFile(caPubKeyPem, gentlsCaPubKeyPath); err != nil {
				lib.Crash(err)
			}
			if err := lib.TlsWriteFile(caCert, gentlsCaCertPath); err != nil {
				lib.Crash(err)
			}
			if err := lib.TlsWriteFile(caPrivKeyPem, gentlsCaPrivKeyPath); err != nil {
				lib.Crash(err)
			}
		} else if gentlsGenCert {
			caPrivKeyPem, err := lib.TlsReadFile(gentlsCaPrivKeyPath)
			if err != nil {
				lib.Crash(err)
			}
			caPrivKey, err := lib.TlsDecodePrivKey(caPrivKeyPem)
			if err != nil {
				lib.Crash(err)
			}
			pubKeyPem, err := lib.TlsReadFile(gentlsPubKeyPath)
			if err != nil {
				lib.Crash(err)
			}
			pubKey, err := lib.TlsDecodePubKey(pubKeyPem)
			if err != nil {
				lib.Crash(err)
			}
			cert, err := lib.TlsGenCert(lib.TlsCertTemplate, pubKey, caPrivKey)
			if err != nil {
				lib.Crash(err)
			}
			if err := lib.TlsWriteFile(cert, gentlsCertPath); err != nil {
				lib.Crash(err)
			}
		} else {
			pubKey, privKey, err := lib.TlsGenKey()
			if err != nil {
				lib.Crash(err)
			}
			pubKeyPem, err := lib.TlsEncodePubKey(pubKey)
			if err != nil {
				lib.Crash(err)
			}
			privKeyPem, err := lib.TlsEncodePrivKey(privKey)
			if err != nil {
				lib.Crash(err)
			}
			if err := lib.TlsWriteFile(pubKeyPem, gentlsPubKeyPath); err != nil {
				lib.Crash(err)
			}
			if err := lib.TlsWriteFile(privKeyPem, gentlsPrivKeyPath); err != nil {
				lib.Crash(err)
			}
		}
	},
}
