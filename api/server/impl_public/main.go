package impl_public

import (
	"context"
	"crypto/ed25519"

	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerPublic struct {
	nf6.UnimplementedNf6PublicServer
	db               *pgxpool.Pool
	tlsCaCert        string
	tlsCaPrivKeyPath string
}

func NewServerPublic(db *pgxpool.Pool, tlsCaCert string, tlsCaPrivKeyPath string) *ServerPublic {
	return &ServerPublic{db: db, tlsCaCert: tlsCaCert, tlsCaPrivKeyPath: tlsCaPrivKeyPath}
}

func (s *ServerPublic) GetCaCert(_ context.Context, in *nf6.None) (*nf6.GetCaCert_Reply, error) {
	return &nf6.GetCaCert_Reply{CaCert: s.tlsCaCert}, nil
}

func (s *ServerPublic) CreateAccount(ctx context.Context, in *nf6.CreateAccount_Request) (*nf6.CreateAccount_Reply, error) {
	if err := lib.ValidateEmail(in.GetEmail()); err != nil {
		return nil, err
	}
	if in.GetSshPubKey() == "" {
		return nil, status.Error(codes.InvalidArgument, "SSH public key must not be empty")
	}
	if in.GetTlsPubKey() == "" {
		return nil, status.Error(codes.InvalidArgument, "TLS public key must not be empty")
	}
	if err := lib.DbCheckNotExists(ctx, s.db, "account", "email", in.GetEmail()); err != nil {
		return nil, err
	}
	if err := lib.DbCheckNotExists(ctx, s.db, "account", "ssh_pub_key", in.GetSshPubKey()); err != nil {
		return nil, err
	}
	if err := lib.DbCheckNotExists(ctx, s.db, "account", "tls_pub_key", in.GetTlsPubKey()); err != nil {
		return nil, err
	}
	cert, err := lib.TlsGenCertUsingPrivKeyFile(lib.TlsCertTemplate, ed25519.PublicKey(in.GetTlsPubKey()), s.tlsCaPrivKeyPath)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to generate a cert using the provided TLS public key")
	}
	query := "insert into account (email, ssh_pub_key, tls_pub_key) values (@email, @ssh_pub_key, @tls_pub_key)"
	args := pgx.NamedArgs{
		"email":       in.GetEmail(),
		"ssh_pub_key": in.GetSshPubKey(),
		"tls_pub_key": in.GetTlsPubKey(),
	}
	if _, err := s.db.Exec(ctx, query, args); err != nil {
		return nil, status.Error(codes.Unknown, "account creation failed")
	}
	return &nf6.CreateAccount_Reply{Cert: string(cert)}, nil
}
