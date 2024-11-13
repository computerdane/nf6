package main

import (
	"context"
	"crypto/ed25519"

	"github.com/computerdane/nf6/nf6"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type Server struct {
	nf6.UnimplementedNf6Server
	db *pgxpool.Pool
}

func (s *Server) WhoAmI(ctx context.Context, in *nf6.WhoAmIRequest) (*nf6.WhoAmIReply, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing peer")
	}

	authInfo, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing transport credentials")
	}

	if len(authInfo.State.VerifiedChains) == 0 || len(authInfo.State.VerifiedChains[0]) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing verified certificate")
	}

	pubkey := string(authInfo.State.VerifiedChains[0][0].PublicKey.(ed25519.PublicKey))

	reply := &nf6.WhoAmIReply{Email: "test", SshPublicKey: "test", SslPublicKey: pubkey}
	return reply, nil
	// err := s.db.QueryRow("select email, ssh_public_key, ssl_public_key from account where $1")
}

func (s *Server) GetMachine(_ context.Context, in *nf6.GetMachineRequest) (*nf6.GetMachineReply, error) {
	return &nf6.GetMachineReply{Address: "fishtank.nf6.sh", JumpAddress: config.domain}, nil
}
