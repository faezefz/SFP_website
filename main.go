package main

import (
	"context"
	"log"

	"github.com/faezefz/SFP_website/api"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource = "postgres://root:secret@localhost:5432/sfp_db?sslmode=disable"
)

func main() {
	dbPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Cannot ping db: %v", err)
	}

	server := api.NewServer(dbPool)

	if err := server.Run(":8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
