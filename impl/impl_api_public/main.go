package impl_api_public

import (
	"context"
	"fmt"
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
	VipWgEndpoint     string
	VipWgPubKey       string
}

func (s *ServerPublic) GetCaCert(_ context.Context, in *nf6.None) (*nf6.GetCaCert_Reply, error) {
	return &nf6.GetCaCert_Reply{CaCert: s.TlsCaCert}, nil
}

func (s *ServerPublic) GetIpv6Info(_ context.Context, in *nf6.None) (*nf6.GetIpv6Info_Reply, error) {
	return &nf6.GetIpv6Info_Reply{GlobalPrefix6: s.IpNet6.String(), AccountPrefix6Len: int32(s.AccountPrefix6Len), VipWgEndpoint: s.VipWgEndpoint, VipWgPubKey: s.VipWgPubKey}, nil
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
	var prefix6 *net.IPNet
	if in.GetPrefix6() == "" {
		randPrefix6, err := lib.RandomIpv6Prefix(s.IpNet6, s.AccountPrefix6Len)
		if err != nil {
			return nil, err
		}
		if err := lib.DbCheckNotExists(ctx, s.Db, "account", "prefix6", randPrefix6); err != nil {
			// TODO: make this less cringe
			return nil, status.Error(codes.AlreadyExists, "somehow we generated an IPv6 prefix for you that is already taken. buy a lottery ticket!")
		}
		prefix6 = randPrefix6
	} else {
		var ip net.IP
		ip, prefix6, err = net.ParseCIDR(in.GetPrefix6())
		if err != nil || ip.To4() != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid prefix6")
		}
	}
	ones, _ := prefix6.Mask.Size()
	if ones != s.AccountPrefix6Len {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("prefix6 must be a /%d", s.AccountPrefix6Len))
	}
	if err := lib.EnsureIpv6PrefixContainsAddr(s.IpNet6, prefix6.IP); err != nil {
		return nil, err
	}
	if err := lib.DbCheckNotExists(ctx, s.Db, "account", "prefix6", prefix6); err != nil {
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
		lib.Warn("account creation failed: ", err)
		return nil, status.Error(codes.Unknown, "account creation failed")
	}
	return &nf6.CreateAccount_Reply{Cert: string(cert)}, nil
}
