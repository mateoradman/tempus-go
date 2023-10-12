package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/mateoradman/tempus/internal/config"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := config.LoadConfig("../../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := pgx.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	testQueries = New(conn)
	exitCode := m.Run()
	os.Exit(exitCode)
}
