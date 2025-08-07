package main

import (
	"context"
	"log"
	_ "net/http"

	"github.com/faezefz/SFP_website/api"
	_ "github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource = "postgres://root:secret@localhost:5432/sfp_db?sslmode=disable"
)

func main() {
	// اتصال به دیتابیس با pgxpool
	dbPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()

	// بررسی اتصال به دیتابیس
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Cannot ping db: %v", err)
	}

	// ایجاد سرور
	server := api.NewServer(dbPool, "secret")

	// راه‌اندازی سرور
	if err := server.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
