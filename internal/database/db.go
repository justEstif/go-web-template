package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *Queries

func InitDB() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return err
	}

	// Verify connection
	if err := pool.Ping(context.Background()); err != nil {
		return err
	}

	DB = New(pool)
	log.Println("Database connection established")
	return nil
}
