package server_secure

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

type ServerSecure struct {
	nf6.UnimplementedNf6SecureServer

	db *pgxpool.Pool
}

func NewServer(db *pgxpool.Pool) *ServerSecure {
	return &ServerSecure{db: db}
}

func (s *ServerSecure) Authenticate(ctx context.Context) (uint64, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return 0, status.Error(codes.Unauthenticated, "missing peer")
	}
	authInfo, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return 0, status.Error(codes.Unauthenticated, "missing transport credentials")
	}
	if len(authInfo.State.VerifiedChains) == 0 || len(authInfo.State.VerifiedChains[0]) == 0 {
		return 0, status.Error(codes.Unauthenticated, "missing verified certificate")
	}

	pubKeyBytes := authInfo.State.VerifiedChains[0][0].PublicKey.(ed25519.PublicKey)
	pubKeyMarshalled, err := x509.MarshalPKIXPublicKey(pubKeyBytes)
	if err != nil {
		log.Fatal(err)
	}
	pubKeyPem := new(bytes.Buffer)
	if err := pem.Encode(pubKeyPem, &pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyMarshalled}); err != nil {
		return 0, status.Error(codes.Unauthenticated, "failed to encode public key")
	}
	pubKey := string(pubKeyPem.Bytes())

	var accountId uint64 = 0
	err = s.db.QueryRow(ctx, "select id from account where ssl_public_key = $1", pubKey).Scan(&accountId)
	if err != nil {
		return 0, status.Error(codes.Unauthenticated, "account does not exist")
	}
	if accountId == 0 {
		return 0, status.Error(codes.Unauthenticated, "account does not exist")
	}

	return accountId, nil
}
