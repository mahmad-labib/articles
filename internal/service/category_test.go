package service

import (
	"context"
	"log"
	"testing"

	"github.com/pressly/goose/v3"
	"main.go/internal/domain"
	"main.go/internal/infra/db/repo"
	"main.go/internal/service/validation"
)

func Test_Category(t *testing.T) {
	adminDB, test := CreateTestDB(t)

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
	validation := validation.New()
	service := NewUserService(repository, validation)

	ctx := context.Background()

	category := domain.Category{
		Name: "Category",
		Slug: "category",
	}

	id, err := service.CreateCategory(ctx, &category)
	if err != nil {
		t.Fatalf("failed to create category with err: %v", err)
	}
	if id == 0 {
		t.Fatalf("expected new category got nil")
	}

	item, error := service.GetCategoryByID(ctx, id)
	if error != nil {
		t.Fatalf("failed to get category with err: %v", error)
	}
	if item == nil {
		t.Fatalf("expected a category got nil")
	}

	list, err := service.ListCategories(ctx, 1, 10)
	if err != nil {
		t.Fatalf("expected no error go %v", err)
	}
	if list == nil {
		t.Fatalf("expected a list of categories got nil")
	}

	err = service.DeleteCategory(ctx, id)
	if err != nil {
		t.Fatalf("expected no error got err: %v", err)
	}

	t.Cleanup(func() {
		testDB.Close()

		_, err := adminDB.Exec(`DROP DATABASE if exists test_db_service`)
		if err != nil {
			t.Fatalf("failed to drop test database: %v", err)
		}

		adminDB.Close()
	})
}

func Test_CategoryInvalid(t *testing.T) {
	adminDB, test := CreateTestDB(t)

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
	validation := validation.New()
	service := NewUserService(repository, validation)

	ctx := context.Background()

	category := domain.Category{
		Name: " ",
		Slug: "category  ",
	}

	id, err := service.CreateCategory(ctx, &category)
	if err == nil {
		t.Fatalf("expected validation error but got none")
	}
	if id != 0 {
		t.Fatalf("expected no category created but got one with ID: %d", id)
	}

	item, error := service.GetCategoryByID(ctx, id)
	if error == nil {
		t.Fatalf("expected validation error but got none")
	}
	if item != nil {
		t.Fatalf("expected no category got %v", item)
	}

	err = service.DeleteCategory(ctx, id)
	if err == nil {
		t.Fatalf("expected validation error but got none")
	}

	t.Cleanup(func() {
		testDB.Close()

		_, err := adminDB.Exec(`DROP DATABASE if exists test_db_service_service`)
		if err != nil {
			t.Fatalf("failed to drop test database: %v", err)
		}

		adminDB.Close()
	})
}
