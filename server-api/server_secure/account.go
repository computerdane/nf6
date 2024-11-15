package server_secure

import (
	"context"

	"github.com/computerdane/nf6/nf6"
	"github.com/jackc/pgx/v5"
)

func (s *ServerSecure) UpdateSshPublicKey(ctx context.Context, in *nf6.UpdateSshPublicKeyRequest) (*nf6.UpdateSshPublicKeyReply, error) {
	accountId, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	query := "update account set ssh_public_key = @ssh_public_key where id = @id"
	args := pgx.NamedArgs{
		"ssh_public_key": in.GetSshPublicKey(),
		"id":             accountId,
	}
	_, err = s.db.Exec(ctx, query, args)
	if err != nil {
		return nil, err
	}

	return &nf6.UpdateSshPublicKeyReply{Success: true}, nil
}
