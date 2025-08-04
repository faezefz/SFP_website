package db

import (
	"context"
	"fmt"

	"github.com/faezefz/SFP_website/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*sqlc.Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db:      db,
		Queries: sqlc.New(db),
	}
}


func (store *Store) execTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return err
	}

	q := sqlc.New(tx)

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
