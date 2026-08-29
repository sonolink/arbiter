package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sonolink/arbiter/internal/config"
	"github.com/sonolink/arbiter/internal/storage"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
	}

	switch os.Args[1] {
	case "migrate":
		runMigrate()
	default:
		fmt.Fprintf(os.Stderr, "arbiter: unknown command %q\n\n", os.Args[1])
		printUsage()
	}
}

func printUsage() {
	fmt.Fprint(os.Stderr, `usage: arbiter <command>

commands:
  migrate: apply pending database migrations
`)
	os.Exit(2)
}

func runMigrate() {
	ctx := context.Background()

	db, err := config.LoadDatabase()
	if err != nil {
		log.Fatalf("arbiter: %v", err)
	}

	store, err := storage.NewStore(ctx, db.URL)
	if err != nil {
		log.Fatalf("arbiter: %v", err)
	}
	defer store.Close()

	if err := store.Migrate(ctx); err != nil {
		log.Fatalf("arbiter: %v", err)
	}

	log.Println("arbiter: migrations applied")
}
