package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v4"
)

var testQueries *Queries

const (
	dbSource = "postgresql://root:secret@localhost:5432/tempus?sslmode=disable"
)

func TestMain(m *testing.M) {
	conn, err := pgx.Connect(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	testQueries = New(conn)
	exitCode := m.Run()
	os.Exit(exitCode)
}
