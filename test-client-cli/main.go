package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	openssl "github.com/computerdane/nf6/lib"
	pb "github.com/computerdane/nf6/nf6"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	baseDir            = flag.String("base-dir", "", "location of api data")
	sslDir             = flag.String("ssl-dir", "", "location of ssl data")
	sshDir             = flag.String("ssh-dir", "", "location of ssh data")
	insecureServerAddr = flag.String("insecure-server-addr", "localhost:6968", "host:port address of insecure api server")
	serverAddr         = flag.String("server-addr", "localhost:6969", "host:port address of secure api server")

	ssl *openssl.Openssl
)

func mkdirAll(dir string) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create directory %s: %v", dir, err)
	}
}

func initSsh() error {
	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-f", "./id_ed25519", "-N", "''", "-q")
	cmd.Dir = *sshDir
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	defer stdin.Close()
	_, err = io.WriteString(stdin, "n\n") // Answer 'n' to prompt to overwrite existing file
	if err != nil {
		return err
	}
	return cmd.Run()
}

func initSsl(insecureClient pb.Nf6InsecureClient) {
	ssl = &openssl.Openssl{Dir: *sslDir}
	err := ssl.GenConfigFile()
	if err != nil {
		log.Fatalf("failed to generate ssl config file: %v", err)
	}
	err = ssl.GenKey("client.key")
	if err != nil {
		log.Fatalf("failed to generate ssl key: %v", err)
	}
	err = ssl.GenCsr("client.key", "client.req")
	if err != nil {
		log.Fatalf("failed to generate ssl csr: %v", err)
	}

	csrBytes, err := os.ReadFile(ssl.GetPath("client.req"))
	if err != nil {
		log.Fatalf("could not read ssl csr: %v", err)
	}

	sshBytes, err := os.ReadFile(*sshDir + "/id_ed25519.pub")
	if err != nil {
		log.Fatalf("could not read ssh pubkey: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	registerReply, err := insecureClient.Register(ctx, &pb.RegisterRequest{Email: "danerieber@gmail.com", SshPublicKey: string(sshBytes), SslCsr: string(csrBytes)})
	if err != nil {
		log.Fatalf("failed to register: %v", err)
	}

	cert := registerReply.GetSslCert()
	err = os.WriteFile(ssl.GetPath("client.crt"), []byte(cert), 0600)
	if err != nil {
		log.Fatalf("failed to write ssl cert: %v", err)
	}
}

func main() {
	flag.Parse()

	if *baseDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("could not find user home dir: %v", err)
		}
		*baseDir = homeDir + "/.nf6/client-cli"
	}
	if *sslDir == "" {
		*sslDir = *baseDir + "/ssl"
	}
	if *sshDir == "" {
		*sshDir = *baseDir + "/ssh"
	}

	mkdirAll(*baseDir)
	mkdirAll(*sslDir)
	mkdirAll(*sshDir)

	insecureConn, err := grpc.NewClient(*insecureServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer insecureConn.Close()

	insecureClient := pb.NewNf6InsecureClient(insecureConn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pingReply, err := insecureClient.Ping(ctx, &pb.PingRequest{Ping: true})

	if err != nil || !pingReply.GetPong() {
		log.Fatalf("failed to ping server: %v", err)
	}

	initSsh()
	initSsl(insecureClient)
}
