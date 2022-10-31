package db

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/mateoradman/tempus/internal/config"
)

type Store interface {
	Querier
	SeedDatabase(ctx context.Context, config config.Config) error
}
type SQLStore struct {
	db *pgx.Conn
	*Queries
}

func NewStore(db *pgx.Conn) Store {
	return SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback(ctx)

	q := New(tx)
	err = fn(q)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (store SQLStore) SeedDatabase(ctx context.Context, config config.Config) error {
	return store.execTx(ctx,
		func(q *Queries) error {
			return seedSuperUser(ctx, q, config)
		},
	)
}
