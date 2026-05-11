// internal/database/postgres.go
package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver — side-effect import

	"github.com/Ubaid-Rza-08/go-rest-api/internal/config"
)

// NewPostgres opens a PostgreSQL connection pool using sqlx.
// It applies connection pool settings from config and pings to verify connectivity.
func NewPostgres(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	// Connection pool tuning
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Verify connectivity with a short deadline
	if err := pingWithRetry(db, 5, 2*time.Second); err != nil {
		return nil, fmt.Errorf("pinging db: %w", err)
	}

	return db, nil
}

// pingWithRetry retries Ping up to maxAttempts times with a delay between each.
func pingWithRetry(db *sqlx.DB, maxAttempts int, delay time.Duration) error {
	var err error
	for i := 1; i <= maxAttempts; i++ {
		if err = db.Ping(); err == nil {
			return nil
		}
		time.Sleep(delay)
	}
	return fmt.Errorf("after %d attempts: %w", maxAttempts, err)
}