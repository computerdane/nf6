package lib

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func DbCheckNotExists(ctx context.Context, db *pgxpool.Pool, table string, column string, value string) error {
	var exists int
	err := db.QueryRow(ctx, fmt.Sprintf("select count(*) from %s where %s = $1", table, column), value).Scan(&exists)
	if err != nil || exists != 0 {
		return status.Error(codes.AlreadyExists, fmt.Sprintf("%s exists with the given %s", table, column))
	}
	return nil
}
