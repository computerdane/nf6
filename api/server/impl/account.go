package impl

import (
	"context"

	"github.com/computerdane/nf6/nf6"
)

func (s *Server) GetAccount(ctx context.Context, in *nf6.None) (*nf6.GetAccount_Reply, error) {
	id, err := s.RequireAccountId(ctx)
	if err != nil {
		return nil, err
	}
	reply := nf6.GetAccount_Reply{}
	if err := s.db.QueryRow(ctx, "select email, ssh_pub_key, tls_pub_key from account where id = $1", id).Scan(&reply.Email, &reply.SshPubKey, &reply.TlsPubKey); err != nil {
		return nil, err
	}
	return &reply, nil
}
