package cmd

import (
	"context"
	"log"
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
	PreRun: requireInsecureClient,
	Run: func(cmd *cobra.Command, args []string) {
		sslPubKeyBytes, err := os.ReadFile(sslPubKeyPath)
		if err != nil {
			log.Fatal(err)
		}

		sshPubKeyBytes, err := os.ReadFile(sshPubKeyPath)
		if err != nil {
			log.Fatal(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		registerReply, err := clientInsecure.Register(ctx, &nf6.RegisterRequest{Email: args[0], SslPublicKey: sslPubKeyBytes, SshPublicKey: sshPubKeyBytes})
		if err != nil {
			log.Fatal(err)
		}

		cert := registerReply.GetSslCert()
		err = os.WriteFile(sslCertPath, cert, 0600)
		if err != nil {
			log.Fatal(err)
		}
	},
}
