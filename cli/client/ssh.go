package client

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:    "ssh [target (host or user@host)]",
	Short:  "SSH into an Nf6 host (defaults to root user if no user given)",
	Args:   cobra.ExactArgs(1),
	PreRun: Connect,
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		user := ""
		name := ""
		if strings.Contains(target, "@") {
			tokens := strings.Split(target, "@")
			if len(tokens) != 2 {
				lib.Crash("Invalid target!")
			}
			user = tokens[0]
			name = tokens[1]
		} else {
			user = "root"
			name = target
		}
		ctx, cancel := lib.Context()
		defer cancel()
		host, err := api.GetHost(ctx, &nf6.GetHost_Request{Name: name})
		if err != nil {
			lib.Crash(err)
		}
		sshPath, err := exec.LookPath("ssh")
		if err != nil {
			lib.Crash(err)
		}
		if err := syscall.Exec(sshPath, []string{"ssh", "-i", sshPrivKeyPath, "-o", "StrictHostKeyChecking=accept-new", fmt.Sprintf("%s@%s", user, host.GetAddr6())}, os.Environ()); err != nil {
			lib.Crash(err)
		}
	},
}
