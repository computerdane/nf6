package impl_public

import (
	"context"

	"github.com/computerdane/nf6/nf6"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServerPublic struct {
	nf6.UnimplementedNf6PublicServer
	db     *pgxpool.Pool
	caCert string
}

func NewServerPublic(db *pgxpool.Pool, caCert string) *ServerPublic {
	return &ServerPublic{db: db, caCert: caCert}
}

func (s *ServerPublic) GetCaCert(_ context.Context, in *nf6.None) (*nf6.GetCaCert_Reply, error) {
	return &nf6.GetCaCert_Reply{CaCert: s.caCert}, nil
}
