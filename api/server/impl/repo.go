package impl

import (
	"context"
	"fmt"

	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateRepo(ctx context.Context, in *nf6.CreateRepo_Request) (*nf6.None, error) {
	accountId, err := s.RequireAccountId(ctx)
	if err != nil {
		return nil, err
	}
	if err := lib.ValidateRepoName(in.GetName()); err != nil {
		return nil, err
	}
	if err := lib.DbCheckNotExists(ctx, s.db, "repo", "name", in.GetName()); err != nil {
		return nil, err
	}
	query := "insert into repo (account_id, name) values (@account_id, @name)"
	args := pgx.NamedArgs{
		"account_id": accountId,
		"name":       in.GetName(),
	}
	if _, err := s.db.Exec(ctx, query, args); err != nil {
		fmt.Println(err)
		return nil, status.Error(codes.Unknown, "repo creation failed")
	}
	return nil, nil
}

func (s *Server) GetRepo(ctx context.Context, in *nf6.GetRepo_Request) (*nf6.GetRepo_Reply, error) {
	accountId, err := s.RequireAccountId(ctx)
	if err != nil {
		return nil, err
	}
	reply := nf6.GetRepo_Reply{}
	query := "select id, name from repo where account_id = @account_id and name = @name"
	args := pgx.NamedArgs{
		"account_id": accountId,
		"name":       in.GetName(),
	}
	if err := s.db.QueryRow(ctx, query, args).Scan(&reply.Id, &reply.Name); err != nil {
		return nil, err
	}
	return &reply, nil
}

func (s *Server) ListRepos(ctx context.Context, in *nf6.None) (*nf6.ListRepos_Reply, error) {
	accountId, err := s.RequireAccountId(ctx)
	if err != nil {
		return nil, err
	}
	reply := nf6.ListRepos_Reply{}
	rows, err := s.db.Query(ctx, "select name from repo where account_id = $1", accountId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to select repos")
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

func (s *Server) UpdateRepo(ctx context.Context, in *nf6.UpdateRepo_Request) (*nf6.None, error) {
	accountId, err := s.RequireAccountId(ctx)
	if err != nil {
		return nil, err
	}
	if in.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "repo id must be non-zero")
	}
	if err := lib.DbCheckAccountOwns(ctx, s.db, "repo", in.GetId(), accountId); err != nil {
		return nil, err
	}
	if in.GetName() != "" {
		if err := lib.ValidateRepoName(in.GetName()); err != nil {
			return nil, err
		}
		if err := lib.DbUpdateUniqueColumnInAccount(ctx, s.db, "repo", "name", in.GetName(), in.GetId(), accountId); err != nil {
			return nil, err
		}
	}
	return nil, nil
}
