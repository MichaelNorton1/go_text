// GetPool returns the global database connection pool, initializing it if necessary.
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

func GetPool(ctx context.Context) (*pgxpool.Pool, error) {
	var err error

	once.Do(func() {
		// Use DATABASE_URL (or fallback to POSTGRES_URL set by Neon Integration)
		connStr := os.Getenv("DATABASE_URL")
		if connStr == "" {
			connStr = os.Getenv("POSTGRES_URL")
		}
		if connStr == "" {
			err = fmt.Errorf("DATABASE_URL environment variable is not set")
			return
		}

		config, parseErr := pgxpool.ParseConfig(connStr)
		if parseErr != nil {
			err = fmt.Errorf("unable to parse connection string: %w", err)
			return
		}

		// --- Neon Serverless + Vercel Settings ---
		// Keep max connections per function low since Vercel scales horizontally
		config.MaxConns = 4
		config.MinConns = 0

		// Close idle connections quickly so Neon compute can scale down to zero when idle
		config.MaxConnIdleTime = 30 * time.Second
		config.MaxConnLifetime = 30 * time.Minute

		p, poolErr := pgxpool.NewWithConfig(ctx, config)
		if poolErr != nil {
			err = fmt.Errorf("unable to create connection pool: %w", poolErr)
			return
		}

		if pingErr := p.Ping(ctx); pingErr != nil {
			p.Close()
			err = fmt.Errorf("unable to ping Neon database: %w", pingErr)
			return
		}

		pool = p
	})

	if pool == nil && err == nil {
		err = fmt.Errorf("database pool failed to initialize")
	}
	if migErr := createTables(ctx, pool); migErr != nil {
		log.Printf("[DB WARNING] Table migration failed: %v", migErr)
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
