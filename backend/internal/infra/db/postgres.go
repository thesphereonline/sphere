package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thesphereonline/sphere/config"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func NewPostgres(cfg *config.Config) (*Postgres, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DB.URL)
	if err != nil {
		return nil, err
	}
	return &Postgres{Pool: pool}, nil
}

func (pg *Postgres) Close() {
	pg.Pool.Close()
}
