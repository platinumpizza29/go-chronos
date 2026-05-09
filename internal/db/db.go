package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	Pool *pgxpool.Pool
)

func Connect(ctx context.Context, dbUrl string) error {
	var err error
	Pool, err = pgxpool.New(ctx, dbUrl)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}
	log.Println("Connected to database")
	return Pool.Ping(ctx)
}

func Close() {
	Pool.Close()
}
