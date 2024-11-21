package lib

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func DbCheckNotExists[T any](ctx context.Context, db *pgxpool.Pool, table string, column string, value T) error {
	var exists int
	err := db.QueryRow(ctx, fmt.Sprintf("select count(*) from %s where %s = $1", table, column), value).Scan(&exists)
	if err != nil || exists != 0 {
		return status.Error(codes.AlreadyExists, fmt.Sprintf("%s exists with the given %s", table, column))
	}
	return nil
}

func DbCheckNotExistsInAccount[T any](ctx context.Context, db *pgxpool.Pool, table string, column string, value T, accountId uint64) error {
	var exists int
	err := db.QueryRow(ctx, fmt.Sprintf("select count(*) from %s where %s = $1 and account_id = $2", table, column), value, accountId).Scan(&exists)
	if err != nil || exists != 0 {
		return status.Error(codes.AlreadyExists, fmt.Sprintf("%s exists with the given %s", table, column))
	}
	return nil
}

func DbCheckAccountOwns(ctx context.Context, db *pgxpool.Pool, table string, id uint64, accountId uint64) error {
	var exists int
	err := db.QueryRow(ctx, fmt.Sprintf("select count(*) from %s where id = $1 and account_id = $2", table), id, accountId).Scan(&exists)
	if err != nil || exists == 0 {
		return status.Error(codes.AlreadyExists, fmt.Sprintf("account cannot access the given %s", table))
	}
	return nil
}

func DbSelectColumn[T any](ctx context.Context, db *pgxpool.Pool, table string, column string, id uint64) (*T, error) {
	var value T
	if err := db.QueryRow(ctx, fmt.Sprintf("select %s from %s where id = $1", column, table), id).Scan(&value); err != nil {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("failed to get %s %s", table, column))
	}
	return &value, nil
}

func DbUpdateColumn[T any](ctx context.Context, db *pgxpool.Pool, table string, column string, value T, id uint64) error {
	if _, err := db.Exec(ctx, fmt.Sprintf("update %s set %s = $1 where id = $2", table, column), value, id); err != nil {
		return status.Error(codes.Unknown, fmt.Sprintf("failed to update %s %s", table, column))
	}
	return nil
}

func DbUpdateUniqueColumn[T comparable](ctx context.Context, db *pgxpool.Pool, table string, column string, value T, id uint64) error {

	existing, err := DbSelectColumn[T](ctx, db, table, column, id)
	if err != nil {
		return err
	}
	if *existing != value {
		if err := DbCheckNotExists(ctx, db, table, column, value); err != nil {
			return err
		}
		if err := DbUpdateColumn(ctx, db, table, column, value, id); err != nil {
			return err
		}
	}
	return nil
}

func DbUpdateUniqueColumnInAccount[T comparable](ctx context.Context, db *pgxpool.Pool, table string, column string, value T, id uint64, accountId uint64) error {

	existing, err := DbSelectColumn[T](ctx, db, table, column, id)
	if err != nil {
		return err
	}
	if *existing != value {
		if err := DbCheckNotExistsInAccount(ctx, db, table, column, value, accountId); err != nil {
			return err
		}
		if err := DbUpdateColumn(ctx, db, table, column, value, id); err != nil {
			return err
		}
	}
	return nil
}
