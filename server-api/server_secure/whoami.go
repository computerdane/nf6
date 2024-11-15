package server_secure

import (
	"context"

	"github.com/computerdane/nf6/nf6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerSecure) WhoAmI(ctx context.Context, in *nf6.WhoAmIRequest) (*nf6.WhoAmIReply, error) {
	accountId, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	reply := &nf6.WhoAmIReply{}

	if err := s.db.QueryRow(ctx, "select email, ssh_public_key, ssl_public_key from account where id = $1", accountId).Scan(&reply.Email, &reply.SshPublicKey, &reply.SslPublicKey); err != nil {
		return nil, status.Error(codes.Unauthenticated, "unknown")
	}

	return reply, nil
}
