package service

import (
	"context"
	"database/sql"
	"log"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"main.go/internal/domain"
	_ "main.go/internal/infra/db/migrations"
	"main.go/internal/infra/db/repo"
)

func createTestDB(t *testing.T) (*sql.DB, *gorm.DB) {
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

	return adminDB, test
}
func Test_User(t *testing.T) {
	adminDB, test := createTestDB(t)

	testDB, err := test.DB()
	if err != nil {
		t.Fatal(err)
	}

	if err := goose.Up(testDB, "../infra/db/migrations"); err != nil {
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

	repository := repo.NewArticlesDB(test)
	service := NewArticleService(repository)

	ctx := context.Background()

	user := domain.User{
		Name:     "user",
		Email:    "user@email.com",
		Password: "12345",
	}

	id, err := service.CreateUser(ctx, &user)
	if err != nil {
		t.Fatalf("failed to create user with err: %v", err)
	}
	if id == 0 {
		t.Fatalf("expected new user got nil")
	}

	item, err := service.GetUserByID(ctx, id)
	if err != nil {
		t.Fatalf("failed to find created user with err: %v", err)
	}
	if item == nil || item.ID == 0 {
		t.Fatalf("expected a user got nil")
	}

	list, err := service.ListUser(ctx, 1, 10)
	if err != nil {
		t.Fatalf("expected no error go %v", err)
	}
	if list == nil {
		t.Fatalf("expected a list of users got nil")
	}

	err = service.DeleteUser(ctx, id)
	if err != nil {
		t.Fatalf("expected no error got err: %v", err)
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

func Test_User_InvalidValues(t *testing.T) {
	adminDB, test := createTestDB(t)

	testDB, err := test.DB()
	if err != nil {
		t.Fatal(err)
	}

	if err := goose.Up(testDB, "../infra/db/migrations"); err != nil {
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

	repository := repo.NewArticlesDB(test)
	service := NewArticleService(repository)

	ctx := context.Background()

	user1 := domain.User{
		Name:     "",
		Email:    "user@email.com",
		Password: "12345",
	}

	id, err := service.CreateUser(ctx, &user1)
	if err == nil {
		t.Fatalf("expected err got nil")
	}
	if id > 0 {
		t.Fatalf("expected no id got %v", id)
	}

	user2 := domain.User{
		Name:     "name",
		Email:    "",
		Password: "12345",
	}

	id, err = service.CreateUser(ctx, &user2)
	if err == nil {
		t.Fatalf("expected err got nil")
	}
	if id > 0 {
		t.Fatalf("expected no id got %v", id)
	}

	user3 := domain.User{
		Name:     "name",
		Email:    "user@email.com",
		Password: "1 2345",
	}

	id, err = service.CreateUser(ctx, &user3)
	if err == nil {
		t.Fatalf("expected err got nil")
	}
	if id > 0 {
		t.Fatalf("expected no id got %v", id)
	}

	item, err := service.GetUserByID(ctx, 1)
	if err == nil {
		t.Fatalf("expected err got nil")
	}
	if item != nil {
		t.Fatalf("expected no user got %v", item)
	}

	list, err := service.ListUser(ctx, 1, 10)
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected no list got %v", list)
	}

	err = service.DeleteUser(ctx, 1)
	if err == nil {
		t.Fatalf("expected an error got nil")
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
