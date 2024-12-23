package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	configPath string
	saveConfig bool

	apiHost          string
	apiPort          int
	apiPortPublic    int
	apiTlsPubKeyPath string
	defaultRepo      string
	output           string
	stateDir         string
	timeout          time.Duration
	tlsCaCertPath    string
	tlsCertPath      string
	tlsPrivKeyPath   string
	tlsPubKeyPath    string
	vipGrpcPort      int
	vipWgPort        int
	wgDeviceName     string
	wgPrivKeyPath    string
	apiTlsPubKey     string

	sshDir string
	tlsDir string

	sshName        string
	sshPrivKeyPath string
	sshPubKeyPath  string

	tlsName   string
	tlsCaName string

	conn       *grpc.ClientConn
	connPublic *grpc.ClientConn

	api       nf6.Nf6Client
	apiPublic nf6.Nf6PublicClient
)

func Init(cmd *cobra.Command) {
	cobra.OnInitialize(InitConfig, InitState)

	cmd.PersistentFlags().StringVar(&configPath, "config-path", "", "path to config file")
	cmd.PersistentFlags().BoolVar(&saveConfig, "save-config", false, "save the flags for this execution to the config file")

	lib.AddOption(cmd, &lib.Option{P: &apiHost, Name: "api-host", Shorthand: "", Value: "localhost", Usage: "server host without port"})
	lib.AddOption(cmd, &lib.Option{P: &apiPort, Name: "api-port", Shorthand: "", Value: 6969, Usage: "server port"})
	lib.AddOption(cmd, &lib.Option{P: &apiPortPublic, Name: "api-port-public", Shorthand: "", Value: 6968, Usage: "server public port"})
	lib.AddOption(cmd, &lib.Option{P: &apiTlsPubKeyPath, Name: "api-tls-pub-key-path", Shorthand: "", Value: "", Usage: "path to the API's TLS public key (for wgserver)"})
	lib.AddOption(cmd, &lib.Option{P: &defaultRepo, Name: "default-repo", Shorthand: "", Value: "main", Usage: "default repo to use for all commands"})
	lib.AddOption(cmd, &lib.Option{P: &output, Name: "output", Shorthand: "", Value: "table", Usage: "output type, json/table"})
	lib.AddOption(cmd, &lib.Option{P: &stateDir, Name: "state-dir", Shorthand: "", Value: "", Usage: "path to state directory"})
	lib.AddOption(cmd, &lib.Option{P: &timeout, Name: "timeout", Shorthand: "", Value: 5 * time.Second, Usage: "timeout for gRPC requests"})
	lib.AddOption(cmd, &lib.Option{P: &tlsCaCertPath, Name: "tls-ca-cert-path", Shorthand: "", Value: "", Usage: "path to TLS ca cert"})
	lib.AddOption(cmd, &lib.Option{P: &tlsCertPath, Name: "tls-cert-path", Shorthand: "", Value: "", Usage: "path to TLS cert"})
	lib.AddOption(cmd, &lib.Option{P: &tlsPrivKeyPath, Name: "tls-priv-key-path", Shorthand: "", Value: "", Usage: "path to TLS private key"})
	lib.AddOption(cmd, &lib.Option{P: &tlsPubKeyPath, Name: "tls-pub-key-path", Shorthand: "", Value: "", Usage: "path to TLS public key"})
	lib.AddOption(cmd, &lib.Option{P: &vipGrpcPort, Name: "vip-grpc-port", Shorthand: "", Value: 6970, Usage: "VIP gRPC port"})
	lib.AddOption(cmd, &lib.Option{P: &vipWgPort, Name: "vip-wg-port", Shorthand: "", Value: 51820, Usage: "VIP WireGuard port"})
	lib.AddOption(cmd, &lib.Option{P: &wgDeviceName, Name: "wg-device-name", Shorthand: "", Value: "wg", Usage: "name of WireGuard interface"})
	lib.AddOption(cmd, &lib.Option{P: &wgPrivKeyPath, Name: "wg-priv-key-path", Shorthand: "", Value: "", Usage: "path to WireGuard private key"})

	cmd.AddCommand(accountCmd)
	cmd.AddCommand(genisoCmd)
	cmd.AddCommand(gensshCmd)
	cmd.AddCommand(gentlsCmd)
	cmd.AddCommand(hostCmd)
	cmd.AddCommand(installCmd)
	cmd.AddCommand(registerCmd)
	cmd.AddCommand(repoCmd)
	cmd.AddCommand(sshCmd)
	cmd.AddCommand(vipCmd)
}

func InitConfig() {
	if configPath == "" {
		if lib.IsDevShell {
			lib.SetHomeConfigPath("dev-nf6")
		} else {
			lib.SetHomeConfigPath("nf6")
		}
	} else {
		lib.SetConfigPath(configPath)
	}
	lib.InitConfig(saveConfig)
	lib.SetTimeout(timeout)
	lib.SetOutputType(output)
}

func InitState() {
	if stateDir == "" {
		if lib.IsDevShell {
			lib.SetHomeStateDir("dev-nf6")
		} else {
			lib.SetHomeStateDir("nf6")
		}
	} else {
		lib.SetStateDir(stateDir)
	}
	lib.AddStateSubDir(&lib.StateSubDir{P: &sshDir, Name: "ssh"})
	lib.AddStateSubDir(&lib.StateSubDir{P: &tlsDir, Name: "tls"})
	lib.InitStateDir()

	sshName = "id_ed25519"
	sshPrivKeyPath = sshDir + "/" + sshName
	sshPubKeyPath = sshDir + "/" + sshName + ".pub"

	tlsName = "client"
	tlsCaName = "ca"
	if tlsPrivKeyPath == "" {
		tlsPrivKeyPath = tlsDir + "/" + tlsName + ".key"
	}
	if tlsPubKeyPath == "" {
		tlsPubKeyPath = tlsDir + "/" + tlsName + ".pub"
	}
	if tlsCertPath == "" {
		tlsCertPath = tlsDir + "/" + tlsName + ".crt"
	}
	if tlsCaCertPath == "" {
		tlsCaCertPath = tlsDir + "/ca.crt"
	}
}

func ConnectPublic(_ *cobra.Command, _ []string) {
	var err error
	connPublic, err = grpc.NewClient(fmt.Sprintf("%s:%d", apiHost, apiPortPublic), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		lib.Crash("failed to connect to public server: ", err)
	}
	apiPublic = nf6.NewNf6PublicClient(connPublic)

	if _, err := os.Stat(tlsCaCertPath); err != nil {
		ctx, cancel := lib.Context()
		defer cancel()
		reply, err := apiPublic.GetCaCert(ctx, nil)
		if err != nil {
			lib.Crash("failed to get server's ca cert: ", err)
		}
		caCert := reply.GetCaCert()
		caCertFile, err := os.OpenFile(tlsCaCertPath, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			lib.Crash("failed to open ca cert file: ", err)
		}
		if _, err := caCertFile.WriteString(caCert); err != nil {
			lib.Crash("failed to write ca cert file: ", err)
		}
	}
}

func Connect(_ *cobra.Command, _ []string) {
	if _, err := os.Stat(tlsCaCertPath); err != nil {
		ConnectPublic(nil, nil)
		connPublic.Close()
	}
	if _, err := os.Stat(tlsCaCertPath); err != nil {
		lib.Crash("ca cert file not found: ", err)
	}
	if _, err := os.Stat(tlsCertPath); err != nil {
		lib.Crash("please register first!")
	}
	if _, err := os.Stat(tlsPrivKeyPath); err != nil {
		lib.Crash("please register first!")
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
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    pool,
		RootCAs:      pool,
	})

	conn, err = grpc.NewClient(fmt.Sprintf("%s:%d", apiHost, apiPort), grpc.WithTransportCredentials(creds), grpc.WithAuthority(lib.TlsName))
	if err != nil {
		lib.Crash("failed to connect to server: ", err)
	}
	api = nf6.NewNf6Client(conn)
	if api == nil {
		lib.Crash("please register first!")
	}
}

func ConnectBoth(cmd *cobra.Command, args []string) {
	ConnectPublic(cmd, args)
	Connect(cmd, args)
}
