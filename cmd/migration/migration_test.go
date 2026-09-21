package main

import (
	"path/filepath"
	"testing"

	"main.go/internal/config"
	"main.go/internal/infra/db"
)

func Test_Migration(t *testing.T) {
	envPath := filepath.Join("..", "..", ".env")
	cfg, err := config.New(envPath)
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}

	db, err := db.NewClient(cfg.DB.DSN)
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}

	if db == nil {
		t.Fatalf("expected db got nil")
	}

	if !db.Migrator().HasTable("users") {
		t.Fatal("expected users table got nil")
	}
	if !db.Migrator().HasTable("articles") {
		t.Fatal("expected articles table got nil")
	}
	if !db.Migrator().HasTable("categories") {
		t.Fatal("expected categories table got nil")
	}
	if !db.Migrator().HasTable("comments") {
		t.Fatal("expected comments table got nil")
	}
	if !db.Migrator().HasTable("tags") {
		t.Fatal("expected tags table got nil")
	}

}
