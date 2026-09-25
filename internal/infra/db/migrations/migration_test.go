package migrations

import (
	"log"
	"testing"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestMigrations(t *testing.T) {
	admin, err := gorm.Open(postgres.Open("postgres://postgres:123@localhost:5432/postgres?sslmode=disable"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	adminDB, err := admin.DB()
	if err != nil {
		t.Fatal(err)
	}
	_, err = adminDB.Exec(`DROP DATABASE test_db`)

	if err := adminDB.Ping(); err != nil {
		t.Fatal(err)
	}

	_, err = adminDB.Exec(`CREATE DATABASE test_db`)
	if err != nil {
		t.Fatal(err)
	}

	test, err := gorm.Open(
		postgres.Open("postgres://postgres:123@localhost:5432/test_db?sslmode=disable"),
		&gorm.Config{},
	)
	if err != nil {
		t.Fatal(err)
	}

	testDB, err := test.DB()
	if err != nil {
		t.Fatal(err)
	}

	if err := goose.Up(testDB, "/"); err != nil {
		log.Fatal(err)
	}

	if !test.Migrator().HasTable("users") {
		t.Fatal("expected users table got nil")
	}
	if !test.Migrator().HasTable("articles") {
		t.Fatal("expected articles table got nil")
	}
	if !test.Migrator().HasTable("categories") {
		t.Fatal("expected categories table got nil")
	}
	if !test.Migrator().HasTable("comments") {
		t.Fatal("expected comments table got nil")
	}
	if !test.Migrator().HasTable("tags") {
		t.Fatal("expected tags table got nil")
	}
	if !test.Migrator().HasTable("article_tags") {
		t.Fatal("expected article_tags table got nil")
	}

	t.Cleanup(func() {
		testDB.Close()

		_, err := adminDB.Exec(`DROP DATABASE test_db`)
		if err != nil {
			t.Fatalf("failed to drop test database: %v", err)
		}

		adminDB.Close()
	})
}
