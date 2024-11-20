package server_secure

import (
	"context"

	"github.com/computerdane/nf6/lib"
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

	if err := lib.ValidateRepoName(in.GetName()); err != nil {
		return nil, err
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

func (s *ServerSecure) ListRepos(ctx context.Context, in *nf6.ListReposRequest) (*nf6.ListReposReply, error) {
	accountId, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(ctx, "select name from repo where account_id = $1", accountId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "no repos found")
	}

	reply := &nf6.ListReposReply{Names: []string{}}

	for rows.Next() {
		var repoName = ""
		err := rows.Scan(&repoName)
		if err != nil {
			return nil, status.Error(codes.Internal, "internal server error")
		}
		reply.Names = append(reply.Names, repoName)
	}

	return reply, nil
}

func (s *ServerSecure) RenameRepo(ctx context.Context, in *nf6.RenameRepoRequest) (*nf6.RenameRepoReply, error) {
	accountId, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	if err := lib.ValidateRepoName(in.GetNewName()); err != nil {
		return nil, err
	}

	query := "select id from repo where account_id = @account_id and name = @name"
	args := pgx.NamedArgs{
		"account_id": accountId,
		"name":       in.GetOldName(),
	}
	repoId := 0
	err = s.db.QueryRow(ctx, query, args).Scan(&repoId)
	if repoId == 0 || err != nil {
		return nil, status.Error(codes.NotFound, "repo not found")
	}

	query = "update repo set name = @name where id = @id"
	args = pgx.NamedArgs{
		"name": in.GetNewName(),
		"id":   repoId,
	}
	if _, err := s.db.Exec(ctx, query, args); err != nil {
		return nil, err
	}

	return &nf6.RenameRepoReply{Success: true}, nil
}
