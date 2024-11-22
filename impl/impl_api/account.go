package impl_api

import (
	"context"

	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
)

func (s *Server) GetAccount(ctx context.Context, in *nf6.None) (*nf6.GetAccount_Reply, error) {
	id, err := s.RequireAccountId(ctx)
	if err != nil {
		return nil, err
	}
	reply := nf6.GetAccount_Reply{}
	if err := s.Db.QueryRow(ctx, "select email, ssh_pub_key, tls_pub_key from account where id = $1", id).Scan(&reply.Email, &reply.SshPubKey, &reply.TlsPubKey); err != nil {
		return nil, err
	}
	return &reply, nil
}

func (s *Server) UpdateAccount(ctx context.Context, in *nf6.UpdateAccount_Request) (*nf6.None, error) {
	id, err := s.RequireAccountId(ctx)
	if err != nil {
		return nil, err
	}
	if in.GetEmail() != "" {
		if err := lib.ValidateEmail(in.GetEmail()); err != nil {
			return nil, err
		}
		if err := lib.DbUpdateUniqueColumn(ctx, s.Db, "account", "email", in.GetEmail(), id); err != nil {
			return nil, err
		}
	}
	if in.GetSshPubKey() != "" {
		if err := lib.DbUpdateUniqueColumn(ctx, s.Db, "account", "ssh_pub_key", in.GetSshPubKey(), id); err != nil {
			return nil, err
		}
	}
	return nil, nil
}
