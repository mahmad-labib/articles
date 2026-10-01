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

func Test_Tag(t *testing.T) {
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

	Tag := domain.Tag{
		Name: "Tag",
	}

	id, err := service.CreateTag(ctx, &Tag)
	if err != nil {
		t.Fatalf("failed to create tag with err: %v", err)
	}
	if id == 0 {
		t.Fatalf("expected new tag got nil")
	}

	id1, err := service.CreateTag(ctx, &Tag)
	if err == nil {
		t.Fatalf("expected error got nil")
	}
	if id1 > 0 {
		t.Fatalf("expected no tag got one")
	}

	item, error := service.GetTagByID(ctx, id)
	if error != nil {
		t.Fatalf("failed to get tag with err: %v", error)
	}
	if item == nil {
		t.Fatalf("expected a tag got nil")
	}

	list, err := service.ListTags(ctx, 1, 10)
	if err != nil {
		t.Fatalf("expected no error go %v", err)
	}
	if list == nil {
		t.Fatalf("expected a list of tags got nil")
	}

	err = service.DeleteTag(ctx, id)
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

func Test_Tag_Invalid(t *testing.T) {
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

	Tag := domain.Tag{
		Name: " ",
	}

	id, err := service.CreateTag(ctx, &Tag)
	if err == nil {
		t.Fatalf("expected error got nil")
	}
	if id > 0 {
		t.Fatalf("expected no tag got one")
	}

	item, error := service.GetTagByID(ctx, id)
	if error == nil {
		t.Fatalf("expected error got nil")
	}
	if item != nil {
		t.Fatalf("expected no tag got one")
	}

	err = service.DeleteTag(ctx, id)
	if err == nil {
		t.Fatalf("expected error got nil")
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
