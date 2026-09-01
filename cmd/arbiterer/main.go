package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/sonolink/arbiterer/internal/config"
	"github.com/sonolink/arbiterer/internal/discord"
	"github.com/sonolink/arbiterer/internal/server"
	"github.com/sonolink/arbiterer/internal/storage"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
	}

	logCfg, err := config.LoadLog()
	if err != nil {
		fmt.Fprintf(os.Stderr, "arbiterer: %v\n", err)
		os.Exit(1)
	}

	slog.SetDefault(slog.New(logCfg.Handler(os.Stderr)))

	switch os.Args[1] {
	case "migrate":
		runMigrate()
	case "serve":
		runServe()
	default:
		fmt.Fprintf(os.Stderr, "arbiterer: unknown command %q\n\n", os.Args[1])
		printUsage()
	}
}

func printUsage() {
	fmt.Fprint(os.Stderr, `usage: arbiterer <command>

commands:
  migrate: apply pending database migrations
`)
	os.Exit(2)
}

func runMigrate() {
	ctx := context.Background()

	pg, err := config.LoadPostgres()
	if err != nil {
		slog.Error("loading configuration", "error", err)
		os.Exit(1)
	}

	store, err := storage.NewStore(ctx, pg.DSN())
	if err != nil {
		slog.Error("connecting to postgres", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	if err := store.Migrate(ctx); err != nil {
		slog.Error("applying migrations", "error", err)
		os.Exit(1)
	}

	slog.Info("migrations applied")
}

func runServe() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("loading configuration", "error", err)
		os.Exit(1)
	}

	store, err := storage.NewStore(ctx, cfg.Postgres.DSN())
	if err != nil {
		slog.Error("connecting to postgres", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	discordClient := discord.NewClient(cfg.Discord)

	srv := server.New(
		cfg.Server,
		slog.Default(),
		store,
		discordClient,
	)
	if err := srv.Run(); err != nil {
		slog.Error("running the server", "error", err)
		os.Exit(1)
	}
}
