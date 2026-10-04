package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/mahmad-labib/articles/internal/config"
	"github.com/mahmad-labib/articles/internal/infra/db"
	"github.com/pressly/goose/v3"
)

func main() {
	envPath := filepath.Join(".env")

	cfg, err := config.New(envPath)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to load config %v", err))

		os.Exit(1)
	}

	conn, err := db.NewClient(cfg.DB.DSN)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := db.Close(conn); err != nil {
			log.Fatal(err)
		}
	}()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	migration, err := db.NewMigration(conn)
	if err != nil {
		log.Fatal(err)
	}

	status := checkMigrationStatus()

	switch status {
	case "rollback":
		if err := migration.Down(); err != nil {
			log.Fatal(err)
		}

	case "reset":
		if err := migration.Reset(); err != nil {
			log.Fatal(err)
		}

	default:
		if err := migration.Up(); err != nil {
			log.Fatal(err)
		}
	}

	slog.Info("Migration successful")
}

func checkMigrationStatus() string {
	if len(os.Args) > 1 {
		if slices.Contains(os.Args[1:], "rollback") {
			return "rollback"
		}

		if slices.Contains(os.Args[1:], "reset") {
			return "reset"
		}
	}

	return "migrate"
}
