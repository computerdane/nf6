package client

import (
	"os"
	"strings"

	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:    "register [email]",
	Short:  "Register a new account",
	Args:   cobra.MaximumNArgs(1),
	PreRun: ConnectPublic,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(sshPrivKeyPath); err != nil {
			gensshCmd.Run(nil, []string{})
		}
		if _, err := os.Stat(tlsPrivKeyPath); err != nil {
			gentlsCmd.Run(nil, []string{})
		}
		email := ""
		if len(args) > 0 {
			email = args[0]
		}
		if err := lib.PromptOrValidate(&email, &promptui.Prompt{
			Label:    "Email",
			Validate: lib.ValidateEmail,
		}); err != nil {
			lib.Crash(err)
		}
		sshPubKey, err := os.ReadFile(sshPubKeyPath)
		if err != nil {
			lib.Crash("failed to read SSH public key: ", err)
		}
		sshPubKeyParts := strings.Split(string(sshPubKey), " ")
		if len(sshPubKeyParts) < 2 {
			lib.Crash("invalid SSH public key")
		}
		sshPubKeyOnly := strings.Join(sshPubKeyParts[:2], " ")
		tlsPubKeyPem, err := lib.TlsReadFile(tlsPubKeyPath)
		if err != nil {
			lib.Crash(err)
		}
		ctx, cancel := lib.Context()
		defer cancel()
		reply, err := apiPublic.CreateAccount(ctx, &nf6.CreateAccount_Request{Email: email, SshPubKey: sshPubKeyOnly, TlsPubKey: string(tlsPubKeyPem)})
		if err != nil {
			lib.Crash(err)
		}
		cert := reply.GetCert()
		if err := lib.TlsWriteFile([]byte(cert), tlsCertPath); err != nil {
			lib.Crash(err)
		}
	},
}
