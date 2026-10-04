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

func Test_Article_create(t *testing.T) {
	adminDB, test := CreateTestDB(t)

	testDB, err := test.DB()
	if err != nil {
		t.Fatal(err)
	}

	if err := goose.Up(testDB, "../infra/db/migrations"); err != nil {
		log.Fatal(err)
	}

	repository := repo.NewArticlesDB(test)
	validation := validation.New()
	service := NewUserService(repository, validation)

	ctx := context.Background()

	user := domain.User{
		Name:     "user",
		Email:    "user@email.com",
		Password: "12345",
	}

	userId, err := service.CreateUser(ctx, &user)
	if err != nil {
		t.Fatalf("failed to create user with err: %v", err)
	}

	if userId == 0 {
		t.Fatalf("expected new user got nil")
	}

	tag := domain.Tag{
		Name: "tag",
	}

	category := domain.Category{
		Name: "category",
		Slug: "category",
	}

	article := domain.Article{
		Title:  "title",
		Slug:   "title",
		Status: "draft",
	}

	categoryId, err := service.CreateCategory(ctx, &category)
	if err != nil {
		t.Fatalf("failed to create category with err: %v", err)
	}

	if categoryId == 0 {
		t.Fatalf("expected new category got nil")
	}

	tagId, err := service.CreateTag(ctx, &tag)
	if err != nil {
		t.Fatalf("failed to create tag with err: %v", err)
	}

	if tagId == 0 {
		t.Fatalf("expected new tag got nil")
	}

	Article, err := service.CreateArticle(ctx, article, userId, categoryId, []domain.Tag{{ID: tagId}})
	if err != nil {
		t.Fatalf("failed to create Article with err: %v", err)
	}

	if Article.ID == 0 {
		t.Fatalf("expected new Article got nil")
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
