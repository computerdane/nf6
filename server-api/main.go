package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

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
	port        = flag.Int("port", 6969, "nf6 api server port")
	sslKeyPath  = flag.String("ssl-key-path", "/var/lib/nf6/server-api/ssl/root.key", "location of nf6 ssl key")
	sslCertPath = flag.String("ssl-cert-path", "/var/lib/nf6/server-api/ssl/root.cert", "location of nf6 ssl cert")
	dbpool      *pgxpool.Pool
	config      Config
)

type Server struct {
	pb.UnimplementedNf6Server
}

func (s *Server) GetMachine(_ context.Context, in *pb.GetMachineRequest) (*pb.GetMachineReply, error) {
	return &pb.GetMachineReply{Address: "fishtank.nf6.sh", JumpAddress: config.domain}, nil
}

func main() {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
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

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		os.Exit(1)
	}

	creds, err := credentials.NewServerTLSFromFile(*sslCertPath, *sslKeyPath)
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
