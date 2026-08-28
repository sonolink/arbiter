package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// github.user_id <=> discord.user_id
type User struct {
	tableName     struct{}  `pg:"users"`
	GitHubUserID  int64     `pg:"github_user_id,pk"`
	DiscordUserID string    `pg:"discord_user_id"`
	CreatedAt     time.Time `pg:"created_at"`
	UpdatedAt     time.Time `pg:"updated_at"`
}

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

func (s *Store) Migrate(ctx context.Context) error {
	return s.db.WithContext(ctx).
		Model((*User)(nil)).
		CreateTable(&orm.CreateTableOptions{IfNotExists: true})
}

func (s *Store) GetByGitHubUserID(ctx context.Context, githubUserID int64) (*User, error) {
	acct := new(User)
	err := s.db.WithContext(ctx).
		Model(acct).
		Where("github_user_id = ?", githubUserID).
		Select()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("storage: select user: %w", err)
	}
	return acct, nil
}

func (s *Store) Upsert(ctx context.Context, githubUserID int64, discordUserID string) (*User, error) {
	var acct User
	err := s.db.WithContext(ctx).RunInTransaction(ctx, func(tx *pg.Tx) error {
		err := tx.Model(&acct).Where("github_user_id = ?", githubUserID).Select()
		if err != nil {
			if errors.Is(err, pg.ErrNoRows) {
				acct = User{
					GitHubUserID:  githubUserID,
					DiscordUserID: discordUserID,
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}
				_, err = tx.Model(&acct).Insert()
				return err
			}
			return err
		}
		acct.DiscordUserID = discordUserID
		acct.UpdatedAt = time.Now()
		_, err = tx.Model(&acct).WherePK().Update()
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("storage: upsert user: %w", err)
	}
	return &acct, nil
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
