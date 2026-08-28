package storage

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
)

type Store struct {
	db *pg.DB
}

func NewStoreFromDSN(dsn string) (*Store, error) {
	opts, err := pg.ParseURL(dsn)
	if err != nil {
		return nil, fmt.Errorf("storage: parse dsn: %w", err)
	}

	if err := ensureDatabase(opts); err != nil {
		return nil, err
	}

	db := pg.Connect(opts)
	if err := db.Ping(context.Background()); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("storage: ping: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func ensureDatabase(opts *pg.Options) error {
	adminOpts := *opts
	adminOpts.Database = "postgres"

	admin := pg.Connect(&adminOpts)
	defer admin.Close()

	if err := admin.Ping(context.Background()); err != nil {
		return fmt.Errorf("storage: admin ping: %w", err)
	}

	var exists bool
	_, err := admin.QueryOne(
		pg.Scan(&exists),
		"SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = ?)",
		opts.Database,
	)
	if err != nil {
		return fmt.Errorf("storage: checking database: %w", err)
	}

	if exists {
		return nil
	}

	if _, err := admin.Exec("CREATE DATABASE " + opts.Database); err != nil {
		return fmt.Errorf("storage: creating database: %w", err)
	}
	return nil
}
