package cmd

import (
	"context"
	"os"

	"github.com/computerdane/nf6/nf6"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(registerCmd)
}

var registerCmd = &cobra.Command{
	Use:    "register [email]",
	Short:  "Register with Nf6",
	Args:   cobra.ExactArgs(1),
	PreRun: RequireInsecureClient,
	Run: func(cmd *cobra.Command, args []string) {
		sslPubKeyBytes, err := os.ReadFile(sslPubKeyPath)
		if err != nil {
			Crash(err)
		}

		sshPubKeyBytes, err := os.ReadFile(sshPubKeyPath)
		if err != nil {
			Crash(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		registerReply, err := clientInsecure.Register(ctx, &nf6.RegisterRequest{Email: args[0], SslPublicKey: sslPubKeyBytes, SshPublicKey: sshPubKeyBytes})
		if err != nil {
			Crash(err)
		}

		cert := registerReply.GetSslCert()
		err = os.WriteFile(sslCertPath, cert, 0600)
		if err != nil {
			Crash(err)
		}
	},
}
