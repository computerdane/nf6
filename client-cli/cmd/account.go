package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/computerdane/nf6/nf6"
	"github.com/spf13/cobra"
)

var (
	newSshPrivKeyPath string
	newSshPubKeyPath  string
)

func init() {
	accountCmd.AddCommand(accountImportSshKeysCmd)

	rootCmd.AddCommand(accountCmd)

	accountImportSshKeysCmd.PersistentFlags().StringVar(&newSshPrivKeyPath, "privateKey", "", "path to SSH private key (default $HOME/.ssh/id_ed25519)")
	accountImportSshKeysCmd.PersistentFlags().StringVar(&newSshPubKeyPath, "publicKey", "", "path to SSH public key (default $HOME/.ssh/id_ed25519.pub)")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		Crash(err)
	}

	newSshPrivKeyPath = homeDir + "/.ssh/id_ed25519"
	newSshPubKeyPath = homeDir + "/.ssh/id_ed25519.pub"
}

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage your account",
}

var accountImportSshKeysCmd = &cobra.Command{
	Use:    "import-ssh-keys",
	Short:  "Import an existing ed25519 keypair",
	PreRun: RequireSecureClient,
	Run: func(cmd *cobra.Command, args []string) {
		// verify that public key is an ed25519 key
		newPrivKey, err := os.ReadFile(newSshPrivKeyPath)
		if err != nil {
			Crash(err)
		}
		newPubKey, err := os.ReadFile(newSshPubKeyPath)
		if err != nil {
			Crash(err)
		}

		// backup existing keys
		backupSuffix := fmt.Sprintf(".%d.bak", time.Now().Unix())
		if err := os.Rename(sshPrivKeyPath, sshPrivKeyPath+backupSuffix); err != nil {
			Crash(err)
		}
		if err := os.Rename(sshPubKeyPath, sshPubKeyPath+backupSuffix); err != nil {
			Crash(err)
		}

		rollBack := func() {
			if err := os.Rename(sshPrivKeyPath+backupSuffix, sshPrivKeyPath); err != nil {
				Crash(err)
			}
			if err := os.Rename(sshPubKeyPath+backupSuffix, sshPubKeyPath); err != nil {
				Crash(err)
			}
		}

		// copy new keys
		if err := os.WriteFile(sshPrivKeyPath, newPrivKey, 0600); err != nil {
			rollBack()
			Crash(err)
		}
		if err := os.WriteFile(sshPubKeyPath, newPubKey, 0644); err != nil {
			rollBack()
			Crash(err)
		}

		// send update to server
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		reply, err := clientSecure.UpdateSshPublicKey(ctx, &nf6.UpdateSshPublicKeyRequest{SshPublicKey: newPubKey})
		if err != nil {
			rollBack()
			Crash(err)
		}
		if !reply.GetSuccess() {
			rollBack()
			Crash()
		}
	},
}
