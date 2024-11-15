package server_secure

import (
	"context"
	"regexp"

	"github.com/computerdane/nf6/nf6"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerSecure) CreateRepo(ctx context.Context, in *nf6.CreateRepoRequest) (*nf6.CreateRepoReply, error) {
	accountId, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	if match, _ := regexp.MatchString(`^[A-Za-z0-9\-_]+$`, in.GetName()); !match {
		return nil, status.Error(codes.InvalidArgument, "repo name must only contain characters A-Z, a-z, 0-9, -, and _")
	}

	query := "select count(*) from repo where account_id = @account_id and name = @name"
	args := pgx.NamedArgs{
		"account_id": accountId,
		"name":       in.GetName(),
	}
	repoExists := 0
	if err := s.db.QueryRow(ctx, query, args).Scan(&repoExists); err != nil {
		return nil, err
	}
	if repoExists != 0 {
		return nil, status.Error(codes.AlreadyExists, "repo already exists")
	}

	query = "insert into repo (account_id, name) values (@account_id, @name)"
	if _, err := s.db.Exec(ctx, query, args); err != nil {
		return nil, err
	}

	return &nf6.CreateRepoReply{Success: true}, nil
}
