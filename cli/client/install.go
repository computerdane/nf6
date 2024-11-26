package client

import (
	_ "embed"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed install.sh
var installSh string

var installCmd = &cobra.Command{
	Use:    "install [device (e.g. /dev/sda)]",
	Short:  "Install NixOS on this machine, pre-configured for use with Nf6",
	Args:   cobra.ExactArgs(1),
	PreRun: Connect,
	Run: func(cmd *cobra.Command, args []string) {
		c := exec.Command("bash", "-s", "-", args[0])
		c.Stdin = strings.NewReader(installSh)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()
	},
}
