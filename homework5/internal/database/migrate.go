package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Prepare(ctx context.Context, dsn string) error {
	db, err := Open(ctx, dsn)
	if err != nil {
		return err
	}
	return db.Close()
}

func Open(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open PostgreSQL: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect to PostgreSQL: %w", err)
	}
	if err := migrate(ctx, db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func migrate(ctx context.Context, db *sql.DB) error {
	files, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	locker, err := lock.NewPostgresSessionLocker(lock.WithLockTimeout(1, 300))
	if err != nil {
		return fmt.Errorf("create migration lock: %w", err)
	}
	provider, err := goose.NewProvider(goose.DialectPostgres, db, files, goose.WithSessionLocker(locker))
	if err != nil {
		return fmt.Errorf("initialize migrations: %w", err)
	}
	if _, err := provider.Up(ctx); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
