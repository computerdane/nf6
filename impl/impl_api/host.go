package impl_api

import (
	"context"
	"fmt"
	"net"

	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateHost(ctx context.Context, in *nf6.CreateHost_Request) (*nf6.None, error) {
	accountId, err := s.RequireAccountId(ctx)
	if err != nil {
		return nil, err
	}
	if err := lib.ValidateHostName(in.GetName()); err != nil {
		return nil, err
	}
	if err := lib.ValidateWireguardKey(in.GetWgPubKey()); err != nil {
		return nil, err
	}
	if err := lib.DbCheckNotExists(ctx, s.Db, "host", "name", in.GetName()); err != nil {
		return nil, err
	}
	if err := lib.DbCheckNotExists(ctx, s.Db, "host", "wg_pub_key", in.GetWgPubKey()); err != nil {
		return nil, err
	}
	prefix6, err := lib.DbSelectColumn[net.IPNet](ctx, s.Db, "account", "prefix6", accountId)
	if err != nil {
		return nil, err
	}
	var addr6 net.IP
	if in.GetAddr6() == "" {
		randAddr6, err := lib.RandomIpv6Addr(prefix6)
		if err != nil {
			return nil, err
		}
		if err := lib.DbCheckNotExists(ctx, s.Db, "host", "addr6", randAddr6.String()); err != nil {
			// TODO: make this less cringe
			return nil, status.Error(codes.AlreadyExists, "somehow we generated an IPv6 address for you that is already taken. buy a lottery ticket!")
		}
		addr6 = randAddr6
	} else {
		addr6 = net.ParseIP(in.GetAddr6())
		if addr6 == nil || addr6.To4() != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid addr6")
		}
	}
	if err := lib.EnsureIpv6PrefixContainsAddr(prefix6, addr6); err != nil {
		return nil, err
	}
	if err := lib.DbCheckNotExists(ctx, s.Db, "host", "addr6", addr6); err != nil {
		return nil, err
	}
	query := "insert into host (account_id, name, addr6, wg_pub_key) values (@account_id, @name, @addr6, @wg_pub_key)"
	args := pgx.NamedArgs{
		"account_id": accountId,
		"name":       in.GetName(),
		"addr6":      addr6.String(),
		"wg_pub_key": in.GetWgPubKey(),
	}
	if _, err := s.Db.Exec(ctx, query, args); err != nil {
		fmt.Println(err)
		return nil, status.Error(codes.Unknown, "host creation failed")
	}

	// create route in WireGuard server
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", s.WgServerGrpcHost, s.WgServerGrpcPort), grpc.WithTransportCredentials(s.Creds), grpc.WithAuthority(lib.TlsName))
	if err != nil {
		lib.Warn("failed to connect to WireGuard server: ", err)
		return nil, status.Error(codes.Internal, "failed to add WireGuard route")
	}
	wgServer := nf6.NewNf6WgClient(conn)
	if wgServer == nil {
		lib.Warn("failed to create WgClient: ", err)
		return nil, status.Error(codes.Internal, "failed to add WireGuard route")
	}
	if _, err := wgServer.CreateRoute(ctx, &nf6.CreateRoute_Request{Addr6: addr6.String(), WgPubKey: in.GetWgPubKey()}); err != nil {
		lib.Warn("failed to create route: ", err)
		return nil, status.Error(codes.Internal, "failed to add WireGuard route")
	}

	return nil, nil
}

func (s *Server) GetHost(ctx context.Context, in *nf6.GetHost_Request) (*nf6.GetHost_Reply, error) {
	accountId, err := s.RequireAccountId(ctx)
	if err != nil {
		return nil, err
	}
	reply := nf6.GetHost_Reply{}
	query := "select id, name, addr6, wg_pub_key, tls_pub_key from host where account_id = @account_id and name = @name"
	args := pgx.NamedArgs{
		"account_id": accountId,
		"name":       in.GetName(),
	}
	var addr6 net.IP
	if err := s.Db.QueryRow(ctx, query, args).Scan(&reply.Id, &reply.Name, &addr6, &reply.WgPubKey, &reply.TlsPubKey); err != nil {
		return nil, err
	}
	reply.Addr6 = addr6.To16().String()
	return &reply, nil
}

func (s *Server) ListHosts(ctx context.Context, in *nf6.None) (*nf6.ListHosts_Reply, error) {
	accountId, err := s.RequireAccountId(ctx)
	if err != nil {
		return nil, err
	}
	reply := nf6.ListHosts_Reply{}
	rows, err := s.Db.Query(ctx, "select name from host where account_id = $1", accountId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to select hosts")
	}
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to scan rows")
		}
		reply.Names = append(reply.Names, name)
	}
	return &reply, nil
}

func (s *Server) UpdateHost(ctx context.Context, in *nf6.UpdateHost_Request) (*nf6.None, error) {
	accountId, err := s.RequireAccountId(ctx)
	if err != nil {
		return nil, err
	}
	if in.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "host id must be non-zero")
	}
	if err := lib.DbCheckAccountOwns(ctx, s.Db, "host", in.GetId(), accountId); err != nil {
		return nil, err
	}
	if in.GetName() != "" {
		if err := lib.ValidateHostName(in.GetName()); err != nil {
			return nil, err
		}
		if err := lib.DbUpdateUniqueColumnInAccount(ctx, s.Db, "host", "name", in.GetName(), in.GetId(), accountId); err != nil {
			return nil, err
		}
	}
	if in.GetAddr6() != "" {
		prefix6, err := lib.DbSelectColumn[net.IPNet](ctx, s.Db, "account", "prefix6", accountId)
		if err != nil {
			return nil, err
		}
		addr6 := net.ParseIP(in.GetAddr6())
		if addr6 == nil || addr6.To4() != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid addr6")
		}
		if err := lib.EnsureIpv6PrefixContainsAddr(prefix6, addr6); err != nil {
			return nil, err
		}
		if err := lib.DbUpdateColumn(ctx, s.Db, "host", "addr6", addr6, in.GetId()); err != nil {
			return nil, err
		}
	}
	if in.GetWgPubKey() != "" {
		if err := lib.ValidateWireguardKey(in.GetWgPubKey()); err != nil {
			return nil, err
		}
		if err := lib.DbUpdateUniqueColumn(ctx, s.Db, "host", "wg_pub_key", in.GetWgPubKey(), in.GetId()); err != nil {
			return nil, err
		}
	}
	if in.GetTlsPubKey() != "" {
		if err := lib.DbUpdateUniqueColumn(ctx, s.Db, "host", "tls_pub_key", in.GetTlsPubKey(), in.GetId()); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (s *Server) WgServer_ListHosts(ctx context.Context, in *nf6.None) (*nf6.WgServer_ListHosts_Reply, error) {
	err := s.RequireWgServerOrigin(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := s.Db.Query(ctx, "select addr6, wg_pub_key from host")
	if err != nil {
		lib.Warn("list hosts query failed: ", err)
		return nil, status.Error(codes.Internal, "failed to get hosts info")
	}
	hosts := []*nf6.WgServer_ListHosts_Reply_Host{}
	for rows.Next() {
		var addr6 net.IP
		var wgPubKey string
		if err := rows.Scan(&addr6, &wgPubKey); err != nil {
			return nil, status.Error(codes.Internal, "failed to scan db data")
		}
		host := &nf6.WgServer_ListHosts_Reply_Host{
			Addr6:    addr6.String(),
			WgPubKey: wgPubKey,
		}
		hosts = append(hosts, host)
	}
	return &nf6.WgServer_ListHosts_Reply{Hosts: hosts}, nil
}
