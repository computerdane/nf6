package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"

	openssl "github.com/computerdane/nf6/lib"
	pb "github.com/computerdane/nf6/nf6"
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

	ssl    *openssl.Openssl
	dbpool *pgxpool.Pool
	config Config
)

type InsecureServer struct {
	pb.UnimplementedNf6InsecureServer
	db *pgxpool.Pool
}

func (s *InsecureServer) Ping(_ context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	if in.GetPing() {
		return &pb.PingResponse{Pong: true}, nil
	}
	return nil, errors.New("did not set ping to true")
}

func (s *InsecureServer) Register(_ context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error) {
	certBytes, err := ssl.GenCertFromCsrInMemory(in.GetSslCsr(), "ca.key", "ca.crt")
	if err != nil {
		log.Printf("failed to generate ssl cert from csr in memory: %v", err)
		return nil, err
	}
	cert := string(certBytes)

	pubkeyBytes, err := ssl.GetPublicKeyInMemory(cert)
	if err != nil {
		log.Printf("failed to get public key from cert: %v", err)
		return nil, err
	}
	pubkey := string(pubkeyBytes)

	_, err = s.db.Exec(context.Background(), "insert into account (email, ssh_public_key, ssl_public_key) values ($1, $2, $3)", in.GetEmail(), in.GetSshPublicKey(), pubkey)
	if err != nil {
		log.Printf("sql query failed: %v", err)
		return nil, err
	}

	return &pb.RegisterReply{SslCert: cert}, nil
}

type Server struct {
	pb.UnimplementedNf6Server
	db *pgxpool.Pool
}

func (s *Server) GetMachine(_ context.Context, in *pb.GetMachineRequest) (*pb.GetMachineReply, error) {
	return &pb.GetMachineReply{Address: "fishtank.nf6.sh", JumpAddress: config.domain}, nil
}

func initSsl() {
	ssl = &openssl.Openssl{Dir: *sslDir}
	err := ssl.GenConfigFile()
	if err != nil {
		log.Fatalf("failed to generate ssl config file: %v", err)
	}
	err = ssl.GenKey("ca.key")
	if err != nil {
		log.Fatalf("failed to generate ssl ca key: %v", err)
	}
	err = ssl.GenCert("ca.key", "ca.crt")
	if err != nil {
		log.Fatalf("failed to generate ssl ca cert: %v", err)
	}
	err = ssl.GenKey("server.key")
	if err != nil {
		log.Fatalf("failed to generate ssl key: %v", err)
	}
	err = ssl.GenCsr("server.key", "server.req")
	if err != nil {
		log.Fatalf("failed to generate ssl csr: %v", err)
	}
	err = ssl.GenCertFromCsr("server.req", "ca.key", "ca.crt", "server.crt")
	if err != nil {
		log.Fatalf("failed to generate ssl cert from csr: %v", err)
	}
}

func main() {
	flag.Parse()

	initSsl()

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

	creds, err := credentials.NewServerTLSFromFile(ssl.GetPath("server.crt"), ssl.GetPath("server.key"))
	if err != nil {
		log.Fatalf("failed to initialize server TLS: %v", err)
	}

	insecureServer := grpc.NewServer()
	server := grpc.NewServer(grpc.Creds(creds))

	pb.RegisterNf6InsecureServer(insecureServer, &InsecureServer{db: dbpool})
	pb.RegisterNf6Server(server, &Server{db: dbpool})

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
