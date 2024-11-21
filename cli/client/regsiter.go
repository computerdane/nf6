package client

import (
	"os"
	"strings"

	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register [email]",
	Short: "Register using your email",
	Args:  cobra.ExactArgs(1),

	PreRun: ConnectPublic,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(sshPrivKeyPath); err != nil {
			gensshCmd.Run(nil, []string{})
		}
		if _, err := os.Stat(tlsPrivKeyPath); err != nil {
			gentlsCmd.Run(nil, []string{})
		}
		email := args[0]
		sshPubKey, err := os.ReadFile(sshPubKeyPath)
		if err != nil {
			lib.Crash("failed to read SSH public key: ", err)
		}
		sshPubKeyParts := strings.Split(string(sshPubKey), " ")
		if len(sshPubKeyParts) < 2 {
			lib.Crash("invalid SSH public key")
		}
		sshPubKeyOnly := strings.Join(sshPubKeyParts[:2], " ")
		tlsPubKey, err := os.ReadFile(tlsPubKeyPath)
		if err != nil {
			lib.Crash("failed to read TLS public key: ", err)
		}
		ctx, cancel := lib.Context()
		defer cancel()
		reply, err := clientPublic.CreateAccount(ctx, &nf6.CreateAccount_Request{Email: email, SshPubKey: sshPubKeyOnly, TlsPubKey: string(tlsPubKey)})
		if err != nil {
			lib.Crash(err)
		}
		cert := reply.GetCert()
		certFile, err := os.OpenFile(tlsCertPath, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			lib.Crash("failed to open cert file: ", err)
		}
		if _, err := certFile.WriteString(cert); err != nil {
			lib.Crash("failed to write cert file: ", err)
		}
	},
}
