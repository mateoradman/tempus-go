package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mateoradman/tempus/internal/config"
)

type Store interface {
	Querier
	SeedDatabase(ctx context.Context, config config.Config) error
}
type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

func NewStore(connPool *pgxpool.Pool) Store {
	return SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
