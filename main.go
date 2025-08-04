package main

import (
	"context"
	"log"

	"github.com/faezefz/SFP_website/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	connString := "postgres://root:secret@localhost:5432/sfp_db?sslmode=disable"

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer pool.Close()

	store := db.NewStore(pool)

}
