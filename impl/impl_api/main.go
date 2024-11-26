package impl_api

import (
	"context"
	"net"

	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type Server struct {
	nf6.UnimplementedNf6Server
	Creds        credentials.TransportCredentials
	Db           *pgxpool.Pool
	IpNet6       *net.IPNet
	VipGrpcHost  string
	VipGrpcPort  int
	VipTlsPubKey string
}

func (s *Server) RequireAccountId(ctx context.Context) (uint64, error) {
	pubKey, err := lib.TlsGetGrpcPubKey(ctx)
	if err != nil {
		return 0, err
	}
	var id uint64 = 0
	if err := s.Db.QueryRow(ctx, "select id from account where tls_pub_key = $1", pubKey).Scan(&id); err != nil || id == 0 {
		return 0, status.Error(codes.Unauthenticated, "account does not exist")
	}
	return id, nil
}

func (s *Server) RequireVipOrigin(ctx context.Context) error {
	pubKey, err := lib.TlsGetGrpcPubKey(ctx)
	if err != nil {
		return err
	}
	if pubKey != s.VipTlsPubKey {
		lib.Warn("attempt made with unknown public key: ", pubKey)
		return status.Error(codes.Unauthenticated, "access denied")
	}
	return nil
}
