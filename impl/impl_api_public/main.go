package impl_api_public

import (
	"context"
	"net"

	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerPublic struct {
	nf6.UnimplementedNf6PublicServer
	AccountPrefix6Len int
	Db                *pgxpool.Pool
	IpNet6            *net.IPNet
	TlsCaCert         string
	TlsCaPrivKeyPath  string
}

func (s *ServerPublic) GetCaCert(_ context.Context, in *nf6.None) (*nf6.GetCaCert_Reply, error) {
	return &nf6.GetCaCert_Reply{CaCert: s.TlsCaCert}, nil
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
	tlsPubKey, err := lib.TlsDecodePubKey([]byte(in.GetTlsPubKey()))
	if err != nil {
		return nil, err
	}
	// TODO: Allow user to set their IPv6 prefix
	// prefix6 := in.GetPrefix6()
	prefix6 := ""
	if prefix6 == "" {
		randPrefix6, err := lib.RandomIpv6Prefix(s.IpNet6, s.AccountPrefix6Len)
		if err != nil {
			return nil, err
		}
		if err := lib.DbCheckNotExists(ctx, s.Db, "account", "prefix6", randPrefix6.String()); err != nil {
			// TODO: make this less cringe
			return nil, status.Error(codes.AlreadyExists, "somehow we generated an IPv6 prefix for you that is already taken. buy a lottery ticket!")
		}
		prefix6 = randPrefix6.String()
	}
	if err := lib.ValidateIpv6Prefix(prefix6, s.AccountPrefix6Len); err != nil {
		return nil, err
	}
	if err := lib.DbCheckNotExists(ctx, s.Db, "account", "email", in.GetEmail()); err != nil {
		return nil, err
	}
	if err := lib.DbCheckNotExists(ctx, s.Db, "account", "ssh_pub_key", in.GetSshPubKey()); err != nil {
		return nil, err
	}
	if err := lib.DbCheckNotExists(ctx, s.Db, "account", "tls_pub_key", in.GetTlsPubKey()); err != nil {
		return nil, err
	}
	if err := lib.DbCheckNotExists(ctx, s.Db, "account", "prefix6", prefix6); err != nil {
		return nil, err
	}
	cert, err := lib.TlsGenCertUsingPrivKeyFile(lib.TlsCertTemplate, tlsPubKey, s.TlsCaPrivKeyPath)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to generate a cert using the provided TLS public key")
	}
	query := "insert into account (email, ssh_pub_key, tls_pub_key, prefix6) values (@email, @ssh_pub_key, @tls_pub_key, @prefix6)"
	args := pgx.NamedArgs{
		"email":       in.GetEmail(),
		"ssh_pub_key": in.GetSshPubKey(),
		"tls_pub_key": in.GetTlsPubKey(),
		"prefix6":     prefix6,
	}
	if _, err := s.Db.Exec(ctx, query, args); err != nil {
		return nil, status.Error(codes.Unknown, "account creation failed")
	}
	return &nf6.CreateAccount_Reply{Cert: string(cert)}, nil
}
