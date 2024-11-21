package client

import (
	"os"
	"os/exec"

	"github.com/computerdane/nf6/lib"
	"github.com/spf13/cobra"
)

var gensshCmd = &cobra.Command{
	Use:   "genssh",
	Short: "Generate a new SSH keypair",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(sshPrivKeyPath); err != nil {
			if err := exec.Command("ssh-keygen", "-t", "ed25519", "-f", sshPrivKeyPath, "-N", "", "-q").Run(); err != nil {
				lib.Crash(err)
			}
		} else {
			lib.Crash("ssh keypair already exists")
		}
	},
}
