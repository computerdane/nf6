package cmd

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/computerdane/nf6/nf6"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	baseDir string
	sslDir  string
	sshDir  string

	privKeyPath string
	pubKeyPath  string
	certPath    string

	serverHost         string
	serverPort         string
	serverPortInsecure string

	timeout time.Duration

	conn         *grpc.ClientConn
	connInsecure *grpc.ClientConn

	client         nf6.Nf6Client
	clientInsecure nf6.Nf6InsecureClient

	rootCmd = &cobra.Command{
		Use:   "nf",
		Short: "nf simplifies OS provisioning and deployment",
	}
)

func mkdirAll(dir string) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create directory %s: %v", dir, err)
	}
}

func requireSecureClient(_ *cobra.Command, _ []string) {
	if client == nil {
		log.Print("error: you must be registered!")
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initDirs, initSsh, initSsl)
	rootCmd.PersistentFlags().StringVar(&baseDir, "base-dir", "", "location of base dir (default ~/.nf6/client-cli)")
	rootCmd.PersistentFlags().StringVar(&sslDir, "ssl-dir", "", "location of ssl dir (default ~/.nf6/client-cli/ssl)")
	rootCmd.PersistentFlags().StringVar(&sshDir, "ssh-dir", "", "location of ssh dir (default ~/.nf6/client-cli/ssh)")
	rootCmd.PersistentFlags().StringVar(&serverHost, "server-host", "localhost", "server host without port")
	rootCmd.PersistentFlags().StringVar(&serverPort, "server-port", "6969", "server secure port")
	rootCmd.PersistentFlags().StringVar(&serverPortInsecure, "server-port-insecure", "6968", "server insecure port")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 10*time.Second, "grpc timeout")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		var err error
		connInsecure, err = grpc.NewClient(serverHost+":"+serverPortInsecure, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

		if _, err := os.Stat(certPath); errors.Is(err, os.ErrNotExist) {
			return
		}

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

		cert, err := tls.LoadX509KeyPair(certPath, privKeyPath)
		if err != nil {
			log.Fatalf("failed to load x509 keypair: %v", err)
		}
		creds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    caCertPool,
			RootCAs:      caCertPool,
		})

		conn, err = grpc.NewClient(serverHost+":"+serverPort, grpc.WithTransportCredentials(creds), grpc.WithAuthority("a"))
		if err != nil {
			log.Fatalf("failed to dial: %v", err)
		}
		client = nf6.NewNf6Client(conn)
	}

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		log.Printf("Thanks for using %s! For help, use `%s help`", cmd.Use, cmd.Use)
	}

	rootCmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {
		if connInsecure != nil {
			connInsecure.Close()
		}
		if conn != nil {
			conn.Close()
		}
	}
}

func initDirs() {
	if baseDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("could not find user home dir: %v", err)
		}
		baseDir = homeDir + "/.nf6/client-cli"
	}
	if sslDir == "" {
		sslDir = baseDir + "/ssl"
	}
	if sshDir == "" {
		sshDir = baseDir + "/ssh"
	}

	mkdirAll(baseDir)
	mkdirAll(sslDir)
	mkdirAll(sshDir)

	privKeyPath = sslDir + "/client.key"
	pubKeyPath = sslDir + "/client.key.pub"
	certPath = sslDir + "/client.crt"
}

func initSsh() {
	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-f", "./id_ed25519", "-N", "''", "-q")
	cmd.Dir = sshDir
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("init ssh: failed to open stdin: %v", err)
	}
	defer stdin.Close()
	_, err = io.WriteString(stdin, "n\n") // Answer 'n' to prompt to overwrite existing file
	if err != nil {
		log.Fatalf("init ssh: failed to wrtie to stdin: %v", err)
	}
	cmd.Run()
}

func initSsl() {
	if _, err := os.Stat(privKeyPath); errors.Is(err, os.ErrNotExist) {
		pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			log.Fatal(err)
		}

		privKeyMarshalled, err := x509.MarshalPKCS8PrivateKey(privKey)
		if err != nil {
			log.Fatal(err)
		}
		privKeyPem, err := os.Create(privKeyPath)
		if err != nil {
			log.Fatal(err)
		}
		pem.Encode(privKeyPem, &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privKeyMarshalled,
		})

		pubKeyMarshalled, err := x509.MarshalPKIXPublicKey(pubKey)
		if err != nil {
			log.Fatal(err)
		}
		pubKeyPem, err := os.Create(pubKeyPath)
		if err != nil {
			log.Fatal(err)
		}
		pem.Encode(pubKeyPem, &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: pubKeyMarshalled,
		})
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
