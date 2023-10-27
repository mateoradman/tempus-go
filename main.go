package main

import (
	"context"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mateoradman/tempus/internal/api"
	"github.com/mateoradman/tempus/internal/config"
	db "github.com/mateoradman/tempus/internal/db/sqlc"
)

func main() {
	// Load environment variables
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Connect to the database
	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	// Run migrations
	runDatabaseMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)
	// Seed the database with pre-defined values
	err = store.SeedDatabase(context.Background(), config)
	if err != nil {
		log.Fatalf("cannot seed database %v", err)
	}

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalf("cannot start server at address %s due to err: %v", config.ServerAddress, err)
	}
}

func runDatabaseMigration(migrationURL string, dbSource string) {
	m, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance:", err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up:", err)
	}
	log.Println("database migration successful!")
}
