package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/mateoradman/tempus/api"
	"github.com/mateoradman/tempus/config"
	db "github.com/mateoradman/tempus/internal/db/sqlc"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := pgx.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalf("cannot start server at address %s due to err: %v", config.ServerAddress, err)
	}
}
