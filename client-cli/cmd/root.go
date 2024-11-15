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
	"io"
	"log"
	"os"
	"os/exec"
	"time"

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
	timeout         time.Duration

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
		fmt.Printf("Thanks for using %s! For help, use `%s help`", cmd.Use, cmd.Use)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		disconnect()
	},
}

func disconnect() {
	if connInsecure != nil {
		connInsecure.Close()
	}
	if connSecure != nil {
		connSecure.Close()
	}
}

func crash(err ...error) {
	if len(err) == 0 {
		color.Red("unknown error!")
	} else {
		color.Red(fmt.Sprintf("%v", err[0]))
	}
	disconnect()
	os.Exit(1)
}

func requireInsecureClient(_ *cobra.Command, _ []string) {
	var err error
	connInsecure, err = grpc.NewClient(apiHost+":"+apiPortInsecure, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	clientInsecure = nf6.NewNf6InsecureClient(connInsecure)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	pingReply, err := clientInsecure.Ping(ctx, &nf6.PingRequest{Ping: true})
	if err != nil || !pingReply.GetPong() {
		log.Fatalf("failed to ping server: %v", err)
	}
}

func requireSecureClient(_ *cobra.Command, _ []string) {
	requireInsecureClient(nil, nil)

	if _, err := os.Stat(sslCertPath); errors.Is(err, os.ErrNotExist) {
		log.Print("error: you must be registered!")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	caCertReply, err := clientInsecure.GetCaCert(ctx, &nf6.GetCaCertRequest{})
	if err != nil {
		log.Fatalf("failed to get ca cert: %v", err)
	}
	caCert := caCertReply.GetCert()

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		log.Fatalf("failed to append ca cert: %v", err)
	}

	cert, err := tls.LoadX509KeyPair(sslCertPath, sslPrivKeyPath)
	if err != nil {
		log.Fatalf("failed to load x509 keypair: %v", err)
	}
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		RootCAs:      caCertPool,
	})

	connSecure, err = grpc.NewClient(apiHost+":"+apiPortSecure, grpc.WithTransportCredentials(creds), grpc.WithAuthority("a"))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	clientSecure = nf6.NewNf6SecureClient(connSecure)

	if clientSecure == nil {
		log.Print("error: you must be registered!")
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initPaths, initSsh, initSsl)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/nf6/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&apiHost, "apiHost", "localhost", "api host without port")
	rootCmd.PersistentFlags().StringVar(&apiPortInsecure, "apiPortInsecure", "6968", "api insecure port")
	rootCmd.PersistentFlags().StringVar(&apiPortSecure, "apiPortSecure", "6969", "api secure port")
	rootCmd.PersistentFlags().StringVar(&dataDir, "dataDir", "", "location of data dir (default is $HOME/.local/share/nf6)")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 10*time.Second, "grpc timeout")

	viper.BindPFlag("apiHost", rootCmd.PersistentFlags().Lookup("apiHost"))
	viper.BindPFlag("apiPortInsecure", rootCmd.PersistentFlags().Lookup("apiPortInsecure"))
	viper.BindPFlag("apiPortSecure", rootCmd.PersistentFlags().Lookup("apiPortSecure"))
	viper.BindPFlag("dataDir", rootCmd.PersistentFlags().Lookup("dataDir"))
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))

}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		viper.AddConfigPath(home + "/.config/nf6")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err == nil {
		apiHost = viper.GetString("apiHost")
		apiPortInsecure = viper.GetString("apiPortInsecure")
		apiPortSecure = viper.GetString("apiPortSecure")
		dataDir = viper.GetString("dataDir")
		timeout = viper.GetDuration("timeout")
	}
}

func initPaths() {
	if dataDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		dataDir = home + "/.local/share/nf6"
	}

	sshDir = dataDir + "/ssh"
	sslDir = dataDir + "/ssl"

	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(sshDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(sslDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	sshPrivKeyPath = sshDir + "/id_ed25519"
	sshPubKeyPath = sshDir + "/id_ed25519.pub"

	sslPrivKeyPath = sslDir + "/client.key"
	sslPubKeyPath = sslDir + "/client.key.pub"
	sslCertPath = sslDir + "/client.crt"

}

func initSsh() {
	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-f", sshPrivKeyPath, "-N", "''", "-q")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stdin.Close()
	_, err = io.WriteString(stdin, "n\n") // Answer 'n' to prompt to overwrite existing file
	if err != nil {
		log.Fatal(err)
	}
	cmd.Run()
}

func initSsl() {
	if _, err := os.Stat(sslPrivKeyPath); errors.Is(err, os.ErrNotExist) {
		pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			log.Fatal(err)
		}

		privKeyMarshalled, err := x509.MarshalPKCS8PrivateKey(privKey)
		if err != nil {
			log.Fatal(err)
		}
		privKeyPem, err := os.OpenFile(sslPrivKeyPath, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatal(err)
		}
		if err := pem.Encode(privKeyPem, &pem.Block{Type: "PRIVATE KEY", Bytes: privKeyMarshalled}); err != nil {
			log.Fatal(err)
		}

		pubKeyMarshalled, err := x509.MarshalPKIXPublicKey(pubKey)
		if err != nil {
			log.Fatal(err)
		}
		pubKeyPem, err := os.OpenFile(sslPubKeyPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if err := pem.Encode(pubKeyPem, &pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyMarshalled}); err != nil {
			log.Fatal(err)
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
