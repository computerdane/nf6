package cmd

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	cfgFile string

	apiHost         string
	apiPortInsecure string
	apiPortSecure   string
	dataDir         string
	defaultRepo     string
	gitHost         string
	timeout         time.Duration

	stringOptions   []lib.StringOption
	durationOptions []lib.DurationOption

	sshDir         string
	sshPrivKeyPath string
	sshPubKeyPath  string

	sslDir         string
	sslCertPath    string
	sslPrivKeyPath string
	sslPubKeyPath  string

	connSecure   *grpc.ClientConn
	connInsecure *grpc.ClientConn

	clientSecure   nf6.Nf6SecureClient
	clientInsecure nf6.Nf6InsecureClient
)

var rootCmd = &cobra.Command{
	Use:   "nf",
	Short: "nf simplifies OS provisioning and deployment",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Thanks for using nf6! For help, use %s\n", color.CyanString("nf help"))
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		Disconnect()
	},
}

func init() {
	cobra.OnInitialize(initConfig, initPaths)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/nf6/config.yaml)")

	stringOptions = []lib.StringOption{
		{P: &apiHost, Name: "apiHost", Value: "localhost", Usage: "api host without port"},
		{P: &apiPortInsecure, Name: "apiPortInsecure", Value: "6968", Usage: "api insecure port"},
		{P: &apiPortSecure, Name: "apiPortSecure", Value: "6969", Usage: "api secure port"},
		{P: &dataDir, Name: "dataDir", Value: "", Usage: "location of data dir (default is $HOME/.local/share/nf6)"},
		{P: &defaultRepo, Name: "defaultRepo", Value: "main", Usage: "default repo to use for all commands"},
		{P: &gitHost, Name: "gitHost", Value: "", Usage: "git host without port (default same as apiHost)"},
	}
	durationOptions = []lib.DurationOption{
		{P: &timeout, Name: "timeout", Value: 10 * time.Second, Usage: "grpc timeout"},
	}

	lib.AddStringOptions(rootCmd, stringOptions)
	lib.AddDurationOptions(rootCmd, durationOptions)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			Crash(err)
		}
		viper.AddConfigPath(home + "/.config/nf6")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err == nil {
		lib.LoadStringOptions(rootCmd, stringOptions)
		lib.LoadDurationOptions(rootCmd, durationOptions)
	}

	// try to generate config file
	cfgFileDir := path.Dir(cfgFile)
	if err := os.MkdirAll(cfgFileDir, os.ModePerm); err != nil {
		Warn("failed to make config directory: ", err)
	}
	if err := viper.WriteConfig(); err != nil {
		Warn("failed to generate config: ", err)
	}

	if gitHost == "" {
		gitHost = apiHost
	}
}

func initPaths() {
	if dataDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			Crash(err)
		}
		dataDir = homeDir + "/.local/share/nf6"
	}

	sshDir = dataDir + "/ssh"
	sslDir = dataDir + "/ssl"

	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		Crash(err)
	}
	if err := os.MkdirAll(sshDir, os.ModePerm); err != nil {
		Crash(err)
	}
	if err := os.MkdirAll(sslDir, os.ModePerm); err != nil {
		Crash(err)
	}

	sshPrivKeyPath = sshDir + "/id_ed25519"
	sshPubKeyPath = sshDir + "/id_ed25519.pub"

	sslPrivKeyPath = sslDir + "/client.key"
	sslPubKeyPath = sslDir + "/client.key.pub"
	sslCertPath = sslDir + "/client.crt"
}

func RequireSsh() {
	if _, err := os.Stat(sshPrivKeyPath); errors.Is(err, os.ErrNotExist) {
		cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-f", sshPrivKeyPath, "-N", "", "-q")
		cmd.Run()
	}
}

func RequireSsl() {
	if _, err := os.Stat(sslPrivKeyPath); errors.Is(err, os.ErrNotExist) {
		pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			Crash(err)
		}

		privKeyMarshalled, err := x509.MarshalPKCS8PrivateKey(privKey)
		if err != nil {
			Crash(err)
		}
		privKeyPem, err := os.OpenFile(sslPrivKeyPath, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			Crash(err)
		}
		if err := pem.Encode(privKeyPem, &pem.Block{Type: "PRIVATE KEY", Bytes: privKeyMarshalled}); err != nil {
			Crash(err)
		}

		pubKeyMarshalled, err := x509.MarshalPKIXPublicKey(pubKey)
		if err != nil {
			Crash(err)
		}
		pubKeyPem, err := os.OpenFile(sslPubKeyPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			Crash(err)
		}
		if err := pem.Encode(pubKeyPem, &pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyMarshalled}); err != nil {
			Crash(err)
		}
	}
}

func RequireInsecureClient(_ *cobra.Command, _ []string) {
	RequireSsh()
	RequireSsl()

	var err error
	connInsecure, err = grpc.NewClient(apiHost+":"+apiPortInsecure, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		Crash("failed to dial: ", err)
	}
	clientInsecure = nf6.NewNf6InsecureClient(connInsecure)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	pingReply, err := clientInsecure.Ping(ctx, &nf6.PingRequest{Ping: true})
	if err != nil || !pingReply.GetPong() {
		Crash("failed to ping server: ", err)
	}
}

func RequireSecureClient(_ *cobra.Command, _ []string) {
	RequireInsecureClient(nil, nil)

	if _, err := os.Stat(sslCertPath); errors.Is(err, os.ErrNotExist) {
		Crash("you must be registered!")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	caCertReply, err := clientInsecure.GetCaCert(ctx, &nf6.GetCaCertRequest{})
	if err != nil {
		Crash("failed to get ca cert: ", err)
	}
	caCert := caCertReply.GetCert()

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		Crash("failed to append ca cert: ", err)
	}

	cert, err := tls.LoadX509KeyPair(sslCertPath, sslPrivKeyPath)
	if err != nil {
		Crash("failed to load x509 keypair: ", err)
	}
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		RootCAs:      caCertPool,
	})

	connSecure, err = grpc.NewClient(apiHost+":"+apiPortSecure, grpc.WithTransportCredentials(creds), grpc.WithAuthority("a"))
	if err != nil {
		Crash("failed to dial: ", err)
	}
	clientSecure = nf6.NewNf6SecureClient(connSecure)

	if clientSecure == nil {
		Crash("you must be registered!")
	}
}

func Disconnect() {
	if connInsecure != nil {
		connInsecure.Close()
	}
	if connSecure != nil {
		connSecure.Close()
	}
}

var (
	red    = color.New(color.FgRed).FprintlnFunc()
	yellow = color.New(color.FgYellow).FprintlnFunc()
)

func Warn(a ...any) {
	yellow(os.Stderr, a...)
}

func Crash(a ...any) {
	if len(a) == 0 {
		red(os.Stderr, "unknown error!")
	} else {
		red(os.Stderr, a...)
	}
	Disconnect()
	os.Exit(1)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		Crash(err)
	}
}
