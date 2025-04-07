package storage

import (
	"context"
	"fmt"
	"time"
	"visitor-counter/config"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	Conn *pgx.Conn
}

func NewStorage(cfg *config.Config) (*Storage, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var conn *pgx.Conn
	var err error

	// Retry logic
	for i := 0; i < 10; i++ { // Retry 10 times
		conn, err = pgx.Connect(context.Background(), connStr)
		if err == nil {
			break
		}
		fmt.Printf("Failed to connect to DB: %v, retrying in 2s...\n", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB after retries: %v", err)
	}

	// Auto-initialize tables
	err = initTables(conn)
	if err != nil {
		return nil, err
	}

	return &Storage{Conn: conn}, nil
}

func initTables(conn *pgx.Conn) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS websites (
            id UUID PRIMARY KEY,
            name TEXT UNIQUE NOT NULL,
            token TEXT UNIQUE NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS pages (
            id UUID PRIMARY KEY,
            website_id UUID REFERENCES websites(id),
            path TEXT NOT NULL,
            visitor_count INT DEFAULT 0,
            UNIQUE (website_id, path)
        )`,
		`CREATE TABLE IF NOT EXISTS visitors (
            ip TEXT NOT NULL,
            page_id UUID REFERENCES pages(id),
            PRIMARY KEY (ip, page_id)
        )`,
	}

	for _, q := range queries {
		_, err := conn.Exec(context.Background(), q)
		if err != nil {
			return err
		}
	}
	return nil
}
