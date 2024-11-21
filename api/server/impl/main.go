package impl

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"

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

func NewServer(db *pgxpool.Pool) *Server {
	return &Server{db: db}
}

func (s *Server) RequireAccountId(ctx context.Context) (uint64, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return 0, status.Error(codes.Unauthenticated, "missing peer info")
	}
	authInfo, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return 0, status.Error(codes.Unauthenticated, "failed to parse TLS info")
	}
	if len(authInfo.State.VerifiedChains) == 0 || len(authInfo.State.VerifiedChains[0]) == 0 {
		return 0, status.Error(codes.Unauthenticated, "missing certificate in chain")
	}
	pubKeyBytes, ok := authInfo.State.VerifiedChains[0][0].PublicKey.(ed25519.PublicKey)
	if !ok {
		return 0, status.Error(codes.Unauthenticated, "failed to parse public key")
	}
	pubKeyMarshalled, err := x509.MarshalPKIXPublicKey(pubKeyBytes)
	if err != nil {
		return 0, status.Error(codes.Unauthenticated, "failed to marshall public key")
	}
	pubKeyPem := new(bytes.Buffer)
	if err := pem.Encode(pubKeyPem, &pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyMarshalled}); err != nil {
		return 0, status.Error(codes.Unauthenticated, "failed to encode public key")
	}
	pubKey := string(pubKeyPem.Bytes())

	var id uint64 = 0
	if err := s.db.QueryRow(ctx, "select id from account where tls_pub_key = $1", pubKey).Scan(&id); err != nil || id == 0 {
		return 0, status.Error(codes.Unauthenticated, "account does not exist")
	}

	return id, nil
}
