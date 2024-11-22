package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"

	"github.com/computerdane/nf6/impl/impl_api"
	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",

	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(tlsCaCertPath); err != nil {
			lib.Crash("ca cert file not found: ", err)
		}
		if _, err := os.Stat(tlsCertPath); err != nil {
			lib.Crash("cert file not found: ", err)
		}
		if _, err := os.Stat(tlsPrivKeyPath); err != nil {
			lib.Crash("private key not found: ", err)
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

		db, err := pgxpool.New(context.Background(), dbUrl)
		if err != nil {
			lib.Crash("failed to connect to db: ", err)
		}
		defer db.Close()

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			lib.Crash("failed to listen: ", err)
		}
		fmt.Printf("listening at %v", lis.Addr())

		server := grpc.NewServer(grpc.Creds(creds))
		nf6.RegisterNf6Server(server, &impl_api.Server{
			Db:     db,
			IpNet6: ipNet6,
		})
		if err := server.Serve(lis); err != nil {
			lib.Crash("failed to serve: ", err)
		}
	},
}
