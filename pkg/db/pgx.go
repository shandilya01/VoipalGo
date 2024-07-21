package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func NewPgxConn(ctx context.Context, databaseUrl string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, databaseUrl)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return conn, nil
}
