package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

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
	baseDir = flag.String("base-dir", "/var/lib/nf6/server-api", "location of api data")
	sslDir  = flag.String("ssl-dir", *baseDir+"/ssl", "location of ssl data")

	port = flag.Int("port", 6969, "nf6 api server port")

	ssl *openssl.Openssl

	dbUrl  = flag.String("db-url", "dbname=nf6", "postgres connection string")
	dbpool *pgxpool.Pool

	config Config
)

type InsecureServer struct {
	pb.UnimplementedNf6InsecureServer
}

// func (s *InsecureServer) Register(_ context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error) {
// 	cert, err := openssl.GenCertWithCsr(in.GetSslCsr(), *sslCaCertPath, *sslCaKeyPath)
// 	if err != nil {
// 		log.Fatalf("failed to gen cert: %v", err)
// 		return nil, err
// 	}

// 	pubkey, err := openssl.PublicKey(cert)
// 	if err != nil {
// 		log.Fatalf("failed to get public key from cert: %v", err)
// 		return nil, err
// 	}

// 	err = dbpool.QueryRow(context.Background(), "insert into account (email, ssh_public_key, ssl_public_key) values ($1, $2, $3)", in.GetEmail(), in.GetSshPublicKey(), pubkey).Scan()

// 	return &pb.RegisterReply{SslCert: cert}, nil
// }

type Server struct {
	pb.UnimplementedNf6Server
}

func (s *Server) GetMachine(_ context.Context, in *pb.GetMachineRequest) (*pb.GetMachineReply, error) {
	return &pb.GetMachineReply{Address: "fishtank.nf6.sh", JumpAddress: config.domain}, nil
}

func initSsl() {
	ssl = &openssl.Openssl{Dir: *sslDir}
	err := ssl.GenConfigFile()
	if err != nil {
		log.Fatalf("failed to generate ssl config file: %v", err)
		os.Exit(1)
	}
	err = ssl.GenKey("ca.key")
	if err != nil {
		log.Fatalf("failed to generate ssl ca key: %v", err)
		os.Exit(1)
	}
	err = ssl.GenCert("ca.key", "ca.crt")
	if err != nil {
		log.Fatalf("failed to generate ssl ca cert: %v", err)
		os.Exit(1)
	}
	err = ssl.GenKey("server.key")
	if err != nil {
		log.Fatalf("failed to generate ssl key: %v", err)
		os.Exit(1)
	}
	err = ssl.GenCsr("server.key", "server.req")
	if err != nil {
		log.Fatalf("failed to generate ssl csr: %v", err)
		os.Exit(1)
	}
	err = ssl.GenCertFromCsr("server.req", "ca.key", "ca.crt", "server.crt")
	if err != nil {
		log.Fatalf("failed to generate ssl cert from csr: %v", err)
		os.Exit(1)
	}
}

func initDb() {
	dbpool, err := pgxpool.New(context.Background(), *dbUrl)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	err = dbpool.QueryRow(context.Background(), "select domain, wireguard_public_key from global_config").Scan(&config.domain, &config.wireguardPublicKey)
	if err != nil {
		log.Fatalf("unable to load global config: %v", err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	initSsl()
	initDb()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		os.Exit(1)
	}

	creds, err := credentials.NewServerTLSFromFile(ssl.GetPath("server.crt"), ssl.GetPath("server.key"))
	if err != nil {
		log.Fatalf("failed to initialize server TLS: %v", err)
		os.Exit(1)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterNf6Server(s, &Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
