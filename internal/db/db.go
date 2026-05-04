package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	migrate "github.com/rubenv/sql-migrate"

	"jaronjones/ace-of-base/internal/config"
)

// Open dials Postgres, pings to verify reachability, applies any pending
// migrations from migrationsFS, and returns the live *sql.DB. The caller is
// responsible for Close().
func Open(ctx context.Context, cfg config.DatabaseConfig, migrationsFS embed.FS, migrationsRoot string) (*sql.DB, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("DATABASE_URL is empty")
	}

	db, err := sql.Open("pgx", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("pinging db: %w", err)
	}

	n, err := migrate.Exec(db, "postgres", &migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrationsFS,
		Root:       migrationsRoot,
	}, migrate.Up)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("running migrations: %w", err)
	}
	slog.Info("db ready", "migrations_applied", n)

	return db, nil
}

// Ping checks the connection is still alive. Used by the /health handler.
func Ping(ctx context.Context, db *sql.DB) error {
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return db.PingContext(pingCtx)
}
