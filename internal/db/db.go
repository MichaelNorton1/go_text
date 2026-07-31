package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
	once sync.Once
)

// GetPool returns the global database connection pool, initializing it if necessary.
func GetPool(ctx context.Context) (*pgxpool.Pool, error) {
	var err error

	once.Do(func() {
		connStr := os.Getenv("DATABASE_URL")
		if connStr == "" {
			err = fmt.Errorf("DATABASE_URL environment variable is missing")
			return
		}

		config, parseErr := pgxpool.ParseConfig(connStr)
		if parseErr != nil {
			err = fmt.Errorf("unable to parse DATABASE_URL: %w", parseErr)
			return
		}

		// Serverless-friendly pool sizing
		config.MaxConns = 5                      // Keep max low so horizontal Vercel scaling doesn't exhaust Postgres
		config.MinConns = 1                      // Minimal warm connections
		config.MaxConnIdleTime = 1 * time.Minute // Drop idle connections quickly

		p, poolErr := pgxpool.NewWithConfig(ctx, config)
		if poolErr != nil {
			err = fmt.Errorf("unable to create connection pool: %w", poolErr)
			return
		}

		if pingErr := p.Ping(ctx); pingErr != nil {
			p.Close()
			err = fmt.Errorf("unable to ping database: %w", pingErr)
			return
		}

		pool = p
		log.Println("[DB] Connection pool initialized successfully!")

		// Run table creation migrations safely on cold boot
		if migErr := createTables(ctx, pool); migErr != nil {
			log.Printf("[DB WARNING] Table migration failed: %v", migErr)
		}
	})

	if pool == nil && err == nil {
		err = fmt.Errorf("database pool failed to initialize")
	}

	return pool, err
}

func createTables(ctx context.Context, p *pgxpool.Pool) error {
	schema := `
	CREATE TABLE IF NOT EXISTS hobbies (
		id SERIAL PRIMARY KEY,
		chat_id BIGINT NOT NULL,
		name VARCHAR(100) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		UNIQUE(chat_id, name)
	);

	CREATE TABLE IF NOT EXISTS hobby_logs (
		id SERIAL PRIMARY KEY,
		chat_id BIGINT NOT NULL,
		hobby_name VARCHAR(100) NOT NULL,
		logged_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);`

	_, err := p.Exec(ctx, schema)
	return err
}
