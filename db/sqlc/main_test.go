package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

const (
	dbSource = "postgres://root:secret@localhost:5432/sfp_db?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := pgx.Connect(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		log.Fatal("cannot ping db:", err)
	}

	testQueries = New(conn)

	exitCode := m.Run()

	conn.Close(context.Background())

	os.Exit(exitCode)
}
