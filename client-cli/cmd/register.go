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
	Use:   "register [email]",
	Short: "Register with Nf6",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sslCsrBytes, err := os.ReadFile(ssl.GetPath("client.req"))
		if err != nil {
			log.Fatalf("could not read ssl csr: %v", err)
		}
		sshPubkeyBytes, err := os.ReadFile(sshDir + "/id_ed25519.pub")
		if err != nil {
			log.Fatalf("could not read ssh pubkey: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		registerReply, err := clientInsecure.Register(ctx, &nf6.RegisterRequest{Email: args[0], SshPublicKey: string(sshPubkeyBytes), SslCsr: string(sslCsrBytes)})
		if err != nil {
			log.Fatalf("failed to register: %v", err)
		}

		cert := registerReply.GetSslCert()
		err = os.WriteFile(ssl.GetPath("client.crt"), []byte(cert), 0600)
		if err != nil {
			log.Fatalf("failed to write ssl cert: %v", err)
		}
	},
}
