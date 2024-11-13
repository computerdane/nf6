package main

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"log"

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

func (s *Server) Authenticate(ctx context.Context) (string, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing peer")
	}

	authInfo, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing transport credentials")
	}

	if len(authInfo.State.VerifiedChains) == 0 || len(authInfo.State.VerifiedChains[0]) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing verified certificate")
	}

	pubKeyBytes := authInfo.State.VerifiedChains[0][0].PublicKey.(ed25519.PublicKey)

	pubKeyMarshalled, err := x509.MarshalPKIXPublicKey(pubKeyBytes)
	if err != nil {
		log.Fatal(err)
	}

	pubKeyPem := new(bytes.Buffer)
	pem.Encode(pubKeyPem, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyMarshalled,
	})
	pubKey := string(pubKeyPem.Bytes())

	var accountExists int
	err = s.db.QueryRow(ctx, "select count(*) from account where ssl_public_key = $1", pubKey).Scan(&accountExists)
	if err != nil {
		return "", err
	}
	if accountExists == 0 {
		return "", status.Error(codes.Unauthenticated, "user does not exist")
	}

	return pubKey, nil
}

func (s *Server) WhoAmI(ctx context.Context, in *nf6.WhoAmIRequest) (*nf6.WhoAmIReply, error) {
	pubKey, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	reply := &nf6.WhoAmIReply{SslPublicKey: pubKey}

	err = s.db.QueryRow(ctx, "select email, ssh_public_key from account where ssl_public_key = $1", pubKey).Scan(&reply.Email, &reply.SshPublicKey)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user not found")
	}

	return reply, nil
}

func (s *Server) GetMachine(_ context.Context, in *nf6.GetMachineRequest) (*nf6.GetMachineReply, error) {
	return &nf6.GetMachineReply{Address: "fishtank.nf6.sh", JumpAddress: config.domain}, nil
}
