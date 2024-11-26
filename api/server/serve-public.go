package server

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/computerdane/nf6/impl/impl_api_public"
	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var servePublicCmd = &cobra.Command{
	Use:   "serve-public",
	Short: "Start the public API server",

	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(tlsCaCertPath); err != nil {
			lib.Crash("ca cert file not found: ", err)
		}

		db, err := pgxpool.New(context.Background(), dbUrl)
		if err != nil {
			lib.Crash("failed to connect to db: ", err)
		}
		defer db.Close()

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", portPublic))
		if err != nil {
			lib.Crash("failed to listen: ", err)
		}
		fmt.Printf("listening at %v", lis.Addr())

		tlsCaCert, err := os.ReadFile(tlsCaCertPath)
		if err != nil {
			lib.Crash("failed to read ca cert file: ", err)
		}

		server := grpc.NewServer()
		nf6.RegisterNf6PublicServer(server, &impl_api_public.ServerPublic{
			AccountPrefix6Len: accountPrefix6Len,
			Db:                db,
			IpNet6:            ipNet6,
			TlsCaCert:         string(tlsCaCert),
			TlsCaPrivKeyPath:  tlsCaPrivKeyPath,
			VipWgEndpoint:     vipWgEndpoint,
			VipWgPubKey:       vipWgPubKey,
		})
		if err := server.Serve(lis); err != nil {
			lib.Crash("failed to serve: ", err)
		}
	},
}
