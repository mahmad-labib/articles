package router

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/jwtauth/v5"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"main.go/internal/domain"
	_ "main.go/internal/infra/db/migrations"
	"main.go/internal/infra/db/repo"
	"main.go/internal/service"
	"main.go/internal/service/validation"
)

func CreateTestDB(t *testing.T) (*sql.DB, *gorm.DB) {
	admin, err := gorm.Open(postgres.Open("postgres://postgres:123@localhost:5432/postgres?sslmode=disable"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	adminDB, err := admin.DB()
	if err != nil {
		t.Fatal(err)
	}

	_, err = adminDB.Exec(`DROP DATABASE if exists test_db_login`)
	if err != nil {
		t.Fatal(err)
	}

	if err := adminDB.Ping(); err != nil {
		t.Fatal(err)
	}

	_, err = adminDB.Exec(`CREATE DATABASE test_db_login`)
	if err != nil {
		t.Fatal(err)
	}

	test, err := gorm.Open(
		postgres.Open("postgres://postgres:123@localhost:5432/test_db_login?sslmode=disable"),
		&gorm.Config{},
	)
	if err != nil {
		t.Fatal(err)
	}

	return adminDB, test
}

func TestFullAuthFlow(t *testing.T) {
	adminDB, test := CreateTestDB(t)

	testDB, err := test.DB()
	if err != nil {
		t.Fatal(err)
	}

	if err := goose.Up(testDB, "../../internal/infra/db/migrations"); err != nil {
		log.Fatal(err)
	}

	repository := repo.NewArticlesDB(test)
	validation := validation.New()
	service := service.NewUserService(repository, validation)

	user := domain.User{
		Name:     "user",
		Email:    "user@example.com",
		Password: "password123",
	}

	ctx := context.Background()

	id, err := service.CreateUser(ctx, &user)
	if err != nil {
		t.Fatalf("failed to create user with err: %v", err)
	}

	if id == 0 {
		t.Fatalf("expected new user got nil")
	}

	var tokenAuth = jwtauth.New(
		"HS256",
		[]byte("asdfasdfsdhfsdiufh"),
		nil,
	)

	router := NewUsersRouter(service, tokenAuth)

	loginBody := bytes.NewBufferString(`{"email":"user@example.com","password":"password123"}`)
	reqLogin := httptest.NewRequest(http.MethodPost, "/users/login", loginBody)
	reqLogin.Header.Set("Content-Type", "application/json")

	recLogin := httptest.NewRecorder()

	router.ServeHTTP(recLogin, reqLogin)

	if recLogin.Code != http.StatusOK {
		t.Fatalf("login failed with status: %d", recLogin.Code)
	}

	var loginResp struct {
		Token string `json:"token"`
	}
	json.Unmarshal(recLogin.Body.Bytes(), &loginResp)

	reqProfile := httptest.NewRequest(http.MethodGet, "/user/profile", nil)
	reqProfile.Header.Set("Authorization", "Bearer "+loginResp.Token)

	recProfile := httptest.NewRecorder()

	router.ServeHTTP(recProfile, reqProfile)

	if recProfile.Code != http.StatusOK {
		t.Errorf("expected status 200 for protected route, got %d", recProfile.Code)
	}

	t.Cleanup(func() {
		testDB.Close()

		_, err := adminDB.Exec(`DROP DATABASE if exists test_db_login`)
		if err != nil {
			t.Fatalf("failed to drop test database: %v", err)
		}

		adminDB.Close()
	})
}
