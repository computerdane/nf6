package main

import (
	"context"

	"github.com/computerdane/nf6/nf6"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerInsecure struct {
	nf6.UnimplementedNf6InsecureServer
	db *pgxpool.Pool
}

func (s *ServerInsecure) Ping(_ context.Context, in *nf6.PingRequest) (*nf6.PingResponse, error) {
	if in.GetPing() {
		return &nf6.PingResponse{Pong: true}, nil
	}
	return nil, status.Error(codes.InvalidArgument, "ping must be true")
}

func (s *ServerInsecure) GetCaCert(_ context.Context, in *nf6.GetCaCertRequest) (*nf6.GetCaCertReply, error) {
	return &nf6.GetCaCertReply{Cert: caCert}, nil
}

func (s *ServerInsecure) Register(ctx context.Context, in *nf6.RegisterRequest) (*nf6.RegisterReply, error) {
	var emailExists int
	err := s.db.QueryRow(ctx, "select count(*) from account where email = $1", in.GetEmail()).Scan(&emailExists)
	if err != nil {
		return nil, err
	}
	if emailExists != 0 {
		return nil, status.Error(codes.AlreadyExists, "user already exists with that email")
	}

	cert, err := ssl.GenCert("ca", in.GetSslPublicKey())
	if err != nil {
		return nil, err
	}

	_, err = s.db.Exec(ctx, "insert into account (email, ssl_public_key, ssh_public_key) values ($1, $2, $3)", in.GetEmail(), in.GetSslPublicKey(), in.GetSshPublicKey())
	if err != nil {
		return nil, err
	}

	return &nf6.RegisterReply{SslCert: cert}, nil
}
