package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/pressly/goose/v3"
	"main.go/internal/config"
	"main.go/internal/infra/db"
	_ "main.go/internal/infra/db/migrations"
)

func main() {
	envPath := filepath.Join(".env")
	cfg, err := config.New(envPath)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to load config %v", err))

		os.Exit(1)
	}

	db, err := db.NewClient(cfg.DB.DSN)
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	status := checkMigrationStatus()
	migrationDir := "internal/infra/db/migrations"

	switch status {
	case "rollback":
		if err := goose.Down(sqlDB, migrationDir); err != nil {
			log.Fatal(err)
		}
		slog.Info("Rollback successful")
	case "reset":
		if err := goose.Reset(sqlDB, migrationDir); err != nil {
			log.Fatal(err)
		}

	default:
		if err := goose.Up(sqlDB, migrationDir); err != nil {
			log.Fatal(err)
		}
		slog.Info("Migration successful")
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
