package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mateoradman/tempus/internal/config"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := config.LoadConfig("../../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	testStore = NewStore(conn)
	exitCode := m.Run()
	os.Exit(exitCode)
}
