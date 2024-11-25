package server

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"
	"time"

	"github.com/computerdane/nf6/lib"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/credentials"
)

var (
	configPath string
	saveConfig bool

	accountPrefix6Len     int
	dbUrl                 string
	globalPrefix6         string
	port                  int
	portPublic            int
	stateDir              string
	timeout               time.Duration
	tlsPrivKeyPath        string
	tlsCertPath           string
	tlsCaPrivKeyPath      string
	tlsCaCertPath         string
	wgServerEndpoint      string
	wgServerGrpcHost      string
	wgServerGrpcPort      int
	wgServerTlsPubKeyPath string
	wgServerWgPubKey      string

	tlsDir    string
	tlsName   string
	tlsCaName string

	ipNet6 *net.IPNet

	creds credentials.TransportCredentials
)

func Init(cmd *cobra.Command) {
	cobra.OnInitialize(InitConfig, InitState)

	cmd.PersistentFlags().StringVar(&configPath, "config-path", "", "path to config file")
	cmd.PersistentFlags().BoolVar(&saveConfig, "save-config", false, "save the flags for this execution to the config file")

	lib.AddOption(cmd, &lib.Option{P: &accountPrefix6Len, Name: "account-prefix6-len", Shorthand: "", Value: 60, Usage: "global ipv6 prefix"})
	lib.AddOption(cmd, &lib.Option{P: &dbUrl, Name: "db-url", Shorthand: "", Value: "dbname=nf6", Usage: "postgres connection string"})
	lib.AddOption(cmd, &lib.Option{P: &globalPrefix6, Name: "global-prefix6", Shorthand: "", Value: "fc69::/48", Usage: "global ipv6 prefix"})
	lib.AddOption(cmd, &lib.Option{P: &port, Name: "port", Shorthand: "", Value: 6969, Usage: "server port"})
	lib.AddOption(cmd, &lib.Option{P: &portPublic, Name: "port-public", Shorthand: "", Value: 6968, Usage: "server public port"})
	lib.AddOption(cmd, &lib.Option{P: &stateDir, Name: "state-dir", Shorthand: "", Value: "", Usage: "path to state directory"})
	lib.AddOption(cmd, &lib.Option{P: &timeout, Name: "timeout", Shorthand: "", Value: 5 * time.Second, Usage: "timeout for gRPC requests"})
	lib.AddOption(cmd, &lib.Option{P: &tlsPrivKeyPath, Name: "tls-priv-key-path", Shorthand: "", Value: "", Usage: "path to this server's TLS private key"})
	lib.AddOption(cmd, &lib.Option{P: &tlsCertPath, Name: "tls-cert-path", Shorthand: "", Value: "", Usage: "path to this server's TLS cert"})
	lib.AddOption(cmd, &lib.Option{P: &tlsCaPrivKeyPath, Name: "tls-ca-priv-key-path", Shorthand: "", Value: "", Usage: "path to this server's TLS ca private key"})
	lib.AddOption(cmd, &lib.Option{P: &tlsCaCertPath, Name: "tls-ca-cert-path", Shorthand: "", Value: "", Usage: "path to the root ca cert"})
	lib.AddOption(cmd, &lib.Option{P: &wgServerEndpoint, Name: "wg-server-endpoint", Shorthand: "", Value: "", Usage: "Endpoint of WireGuard server"})
	lib.AddOption(cmd, &lib.Option{P: &wgServerGrpcHost, Name: "wg-server-grpc-host", Shorthand: "", Value: "localhost", Usage: "WireGuard server host for gRPC"})
	lib.AddOption(cmd, &lib.Option{P: &wgServerGrpcPort, Name: "wg-server-grpc-port", Shorthand: "", Value: 6970, Usage: "WireGuard server port for gRPC"})
	lib.AddOption(cmd, &lib.Option{P: &wgServerTlsPubKeyPath, Name: "wg-server-tls-pub-key-path", Shorthand: "", Value: "", Usage: "path to TLS public key for WireGuard server"})
	lib.AddOption(cmd, &lib.Option{P: &wgServerWgPubKey, Name: "wg-server-wg-pub-key", Shorthand: "", Value: "", Usage: "WireGuard public key for WireGuard server"})

	cmd.AddCommand(serveCmd)
	cmd.AddCommand(servePublicCmd)
}

func InitConfig() {
	if configPath == "" {
		if lib.IsDevShell {
			lib.SetHomeConfigPath("dev-nf6-api")
		} else {
			lib.SetSystemConfigPath("nf6-api")
		}
	} else {
		lib.SetConfigPath(configPath)
	}
	lib.InitConfig(saveConfig)
	lib.SetTimeout(timeout)

	var err error
	_, ipNet6, err = net.ParseCIDR(globalPrefix6)
	if err != nil {
		lib.Crash(err)
	}
	ones, bits := ipNet6.Mask.Size()
	if bits != 128 || ones >= 128 {
		lib.Crash("Invalid global IPv6 prefix")
	}
	if ones >= accountPrefix6Len {
		lib.Crash("The global IPv6 prefix length must be smaller than the account IPv6 prefix length")
	}
	if wgServerEndpoint == "" {
		lib.Crash("You must set the WireGuard server's WireGuard endpoint")
	}
	if wgServerGrpcHost == "" {
		lib.Crash("You must set the WireGuard server's gRPC host")
	}
	if wgServerTlsPubKeyPath == "" {
		lib.Crash("You must set the path to the WireGuard server's TLS public key")
	}
	if wgServerWgPubKey == "" {
		lib.Crash("You must set the WireGuard server's WireGuard public key")
	}
}

func InitState() {
	if stateDir == "" {
		if lib.IsDevShell {
			lib.SetHomeStateDir("dev-nf6-api")
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
	if tlsCaPrivKeyPath == "" {
		tlsCaPrivKeyPath = tlsDir + "/ca.key"
	}
	if tlsCaCertPath == "" {
		tlsCaCertPath = tlsDir + "/ca.crt"
	}

	if _, err := os.Stat(tlsCaCertPath); err != nil {
		lib.Crash("ca cert file not found: ", err)
	}
	if _, err := os.Stat(tlsCertPath); err != nil {
		lib.Crash("cert file not found: ", err)
	}
	if _, err := os.Stat(tlsPrivKeyPath); err != nil {
		lib.Crash("priv key file not found: ", err)
	}

	caCert, err := os.ReadFile(tlsCaCertPath)
	if err != nil {
		lib.Crash("failed to read ca cert: ", err)
	}
	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM(caCert); !ok {
		lib.Crash("failed to append ca cert")
	}
	cert, err := tls.LoadX509KeyPair(tlsCertPath, tlsPrivKeyPath)
	if err != nil {
		lib.Crash("failed to load x509 keypair: ", err)
	}
	creds = credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    pool,
		RootCAs:      pool,
	})
}
