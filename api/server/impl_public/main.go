package impl_public

import (
	"context"

	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerPublic struct {
	nf6.UnimplementedNf6PublicServer
	db        *pgxpool.Pool
	tlsCaCert string
	tlsDir    string
	tlsCaName string
}

func NewServerPublic(db *pgxpool.Pool, tlsCaCert string, tlsDir string, tlsCaName string) *ServerPublic {
	return &ServerPublic{db: db, tlsCaCert: tlsCaCert, tlsDir: tlsDir, tlsCaName: tlsCaName}
}

func (s *ServerPublic) GetCaCert(_ context.Context, in *nf6.None) (*nf6.GetCaCert_Reply, error) {
	return &nf6.GetCaCert_Reply{CaCert: s.tlsCaCert}, nil
}

func (s *ServerPublic) CreateAccount(ctx context.Context, in *nf6.CreateAccount_Request) (*nf6.CreateAccount_Reply, error) {
	if in.GetSshPubKey() == "" {
		return nil, status.Error(codes.InvalidArgument, "SSH public key must not be empty")
	}
	if in.GetTlsPubKey() == "" {
		return nil, status.Error(codes.InvalidArgument, "TLS public key must not be empty")
	}
	if err := lib.DbCheckNotExists(ctx, s.db, "account", "email", in.GetEmail()); err != nil {
		return nil, err
	}
	cert, err := lib.GenCert(s.tlsDir, s.tlsCaName, []byte(in.GetTlsPubKey()))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to generate a cert using the provided TLS public key")
	}
	query := "insert into account (email, ssh_pub_key, tls_pub_key) values (@email, @ssh_pub_key, @tls_pub_key)"
	args := pgx.NamedArgs{
		"email":       in.GetEmail(),
		"ssh_pub_key": in.GetSshPubKey(),
		"tls_pub_key": in.GetTlsPubKey(),
	}
	_, err = s.db.Exec(ctx, query, args)
	if err != nil {
		return nil, status.Error(codes.Unknown, "account creation failed")
	}
	return &nf6.CreateAccount_Reply{Cert: string(cert)}, nil
}
