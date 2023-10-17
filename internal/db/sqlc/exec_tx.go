package db

import (
	"context"
	"fmt"

	"github.com/mateoradman/tempus/internal/config"
)

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
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
