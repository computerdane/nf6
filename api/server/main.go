package server

import (
	"time"

	"github.com/computerdane/nf6/lib"
	"github.com/spf13/cobra"
)

var (
	configPath string
	saveConfig bool

	dbUrl          string
	port           int
	portPublic     int
	stateDir       string
	timeout        time.Duration
	tlsPrivKeyPath string
	tlsCertPath    string
	tlsCaCertPath  string

	tlsDir    string
	tlsName   string
	tlsCaName string
)

func Init(cmd *cobra.Command) {
	cobra.OnInitialize(InitConfig, InitState)

	cmd.PersistentFlags().StringVar(&configPath, "config-path", "", "path to config file")
	cmd.PersistentFlags().BoolVar(&saveConfig, "save-config", false, "save the flags for this execution to the config file")

	lib.AddOption(cmd, &lib.Option{P: &dbUrl, Name: "db-url", Shorthand: "", Value: "dbname=nf6", Usage: "postgres connection string"})
	lib.AddOption(cmd, &lib.Option{P: &port, Name: "port", Shorthand: "", Value: 6969, Usage: "server port"})
	lib.AddOption(cmd, &lib.Option{P: &portPublic, Name: "port-public", Shorthand: "", Value: 6968, Usage: "server public port"})
	lib.AddOption(cmd, &lib.Option{P: &stateDir, Name: "state-dir", Shorthand: "", Value: "", Usage: "path to state directory"})
	lib.AddOption(cmd, &lib.Option{P: &timeout, Name: "timeout", Shorthand: "", Value: 5 * time.Second, Usage: "timeout for gRPC requests"})
	lib.AddOption(cmd, &lib.Option{P: &tlsPrivKeyPath, Name: "tls-private-key-path", Shorthand: "", Value: "", Usage: "path to this server's TLS private key"})
	lib.AddOption(cmd, &lib.Option{P: &tlsCertPath, Name: "tls-cert-path", Shorthand: "", Value: "", Usage: "path to this server's TLS cert"})
	lib.AddOption(cmd, &lib.Option{P: &tlsCaCertPath, Name: "tls-ca-cert-path", Shorthand: "", Value: "", Usage: "path to the root ca cert"})

	cmd.AddCommand(serveCmd)
	cmd.AddCommand(servePublicCmd)
}

func InitConfig() {
	if configPath == "" {
		if lib.IsDevShell {
			lib.SetHomeConfigPath("nf6-api-dev")
		} else {
			lib.SetSystemConfigPath("nf6-api")
		}
	} else {
		lib.SetConfigPath(configPath)
	}
	lib.InitConfig(saveConfig)
	lib.SetTimeout(timeout)
}

func InitState() {
	if stateDir == "" {
		if lib.IsDevShell {
			lib.SetHomeStateDir("nf6-api-dev")
		} else {
			lib.SetSystemStateDir("nf6-api")
		}
	} else {
		lib.SetStateDir(stateDir)
	}
	lib.AddStateSubDir(&lib.StateSubDir{P: &tlsDir, Name: "tls"})
	lib.InitStateDir()

	tlsName = "server"
	tlsCaName = "ca"
	if tlsPrivKeyPath == "" {
		tlsPrivKeyPath = tlsDir + "/" + tlsName + ".key"
	}
	if tlsCertPath == "" {
		tlsCertPath = tlsDir + "/" + tlsName + ".crt"
	}
	if tlsCaCertPath == "" {
		tlsCaCertPath = tlsDir + "/ca.crt"
	}
}
