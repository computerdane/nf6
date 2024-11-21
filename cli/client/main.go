package client

import (
	"time"

	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	configPath string
	saveConfig bool

	defaultRepo string
	host        string
	port        int
	portPublic  int
	stateDir    string
	timeout     time.Duration

	sshDir string
	tlsDir string

	conn       *grpc.ClientConn
	connPublic *grpc.ClientConn

	client       nf6.Nf6Client
	clientPublic nf6.Nf6PublicClient
)

func InitConfig() {
	if configPath == "" {
		lib.SetHomeConfigPath("nf6")
	} else {
		lib.SetConfigPath(configPath)
	}
	lib.InitConfig(saveConfig)
}

func InitState() {
	if stateDir == "" {
		lib.SetHomeStateDir("nf6")
	} else {
		lib.SetStateDir(stateDir)
	}
	lib.AddStateSubDir(&lib.StateSubDir{P: &sshDir, Name: "ssh"})
	lib.AddStateSubDir(&lib.StateSubDir{P: &tlsDir, Name: "tls"})
	lib.InitStateDir()
}

func Init(cmd *cobra.Command) {
	cobra.OnInitialize(InitConfig, InitState)

	cmd.PersistentFlags().StringVar(&configPath, "config-path", "", "path to config file")
	cmd.PersistentFlags().BoolVar(&saveConfig, "save-config", false, "save the flags for this execution to the config file")

	lib.AddOption(cmd, &lib.Option{P: &defaultRepo, Name: "default-repo", Shorthand: "", Value: "main", Usage: "default repo to use for all commands"})
	lib.AddOption(cmd, &lib.Option{P: &host, Name: "host", Shorthand: "H", Value: "nf6.sh", Usage: "server host without port"})
	lib.AddOption(cmd, &lib.Option{P: &port, Name: "port", Shorthand: "", Value: 6969, Usage: "server port"})
	lib.AddOption(cmd, &lib.Option{P: &portPublic, Name: "port-public", Shorthand: "", Value: 6968, Usage: "server public port"})
	lib.AddOption(cmd, &lib.Option{P: &stateDir, Name: "state-dir", Shorthand: "", Value: "", Usage: "path to state directory"})
	lib.AddOption(cmd, &lib.Option{P: &timeout, Name: "timeout", Shorthand: "", Value: 5 * time.Second, Usage: "timeout for gRPC requests"})

	cmd.AddCommand(keygenCmd)
	cmd.AddCommand(registerCmd)
}
