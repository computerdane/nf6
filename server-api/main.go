package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/computerdane/nf6/nf6"
	"github.com/computerdane/nf6/server-api/ssl_util"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Config struct {
	domain             string
	wireguardPublicKey string
}

var (
	insecurePort = flag.Int("insecure-port", 6968, "api server insecure port")
	port         = flag.Int("port", 6969, "api server port")
	baseDir      = flag.String("base-dir", "/var/lib/nf6/server-api", "location of api data")
	sslDir       = flag.String("ssl-dir", *baseDir+"/ssl", "location of ssl data")
	dbUrl        = flag.String("db-url", "dbname=nf6", "postgres connection string")

	ssl    *ssl_util.SslUtil
	caCert []byte
	creds  credentials.TransportCredentials
	dbpool *pgxpool.Pool
	config Config
)

func main() {
	flag.Parse()

	ssl = &ssl_util.SslUtil{Dir: *sslDir}

	err := ssl.GenCaFiles("ca")
	if err != nil {
		log.Fatalf("failed to generate ca files: %v", err)
	}
	err = ssl.GenCertFiles("ca", "server")
	if err != nil {
		log.Fatalf("failed to generate server cert files: %v", err)
	}

	dbpool, err := pgxpool.New(context.Background(), *dbUrl)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}
	defer dbpool.Close()

	err = dbpool.QueryRow(context.Background(), "select domain, wireguard_public_key from global_config").Scan(&config.domain, &config.wireguardPublicKey)
	if err != nil {
		log.Fatalf("unable to load global config: %v", err)
	}

	insecureLis, err := net.Listen("tcp", fmt.Sprintf(":%d", *insecurePort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	caCert, err = os.ReadFile(*sslDir + "/ca.crt")
	if err != nil {
		log.Fatalf("failed to read ca.crt: %v", err)
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		log.Fatalf("failed to append ca cert: %v", err)
	}

	cert, err := tls.LoadX509KeyPair(*sslDir+"/server.crt", *sslDir+"/server.key")
	if err != nil {
		log.Printf("failed to load x509 keypair: %v", err)
		return
	}
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		RootCAs:      caCertPool,
	})

	insecureServer := grpc.NewServer()
	server := grpc.NewServer(grpc.Creds(creds))

	nf6.RegisterNf6InsecureServer(insecureServer, &ServerInsecure{db: dbpool})
	nf6.RegisterNf6Server(server, &Server{db: dbpool})

	go func() {
		log.Printf("server listening at %v", insecureLis.Addr())
		if err := insecureServer.Serve(insecureLis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
