package impl

import (
	"github.com/computerdane/nf6/nf6"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	nf6.UnimplementedNf6Server
	db *pgxpool.Pool
}

func NewServer(db *pgxpool.Pool) *Server {
	return &Server{db: db}
}
