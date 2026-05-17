package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/project-go-sender-recommendation-agent/internal/infrastructure/config"
)

// NewPostgresConnection opens and verifies a PostgreSQL connection.
func NewPostgresConnection(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("database.NewPostgresConnection open: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database.NewPostgresConnection ping: %w", err)
	}

	log.Println("Connected to PostgreSQL database")
	return db, nil
}
