package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/mahmad-labib/articles/internal/config"
	"github.com/mahmad-labib/articles/internal/infra/db"
	"github.com/mahmad-labib/articles/internal/infra/server/admin"
	"github.com/mahmad-labib/articles/internal/infra/server/admin/handlers"
	"github.com/mahmad-labib/articles/internal/service"
	"github.com/mahmad-labib/articles/internal/service/validation"
)

func main() {
	cfg := loadConfig()

	conn, err := db.NewClient(cfg.DB.DSN)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := db.Close(conn); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	repository := db.NewRepository(conn)

	jwt := jwtauth.New(
		"HS256",
		[]byte(cfg.JWT),
		nil,
	)

	userService := loadUserService(repository)
	categoryService := loadCategoryService(repository)
	tagService := loadTagService(repository)
	articleService := loadArticleService(repository)

	handler := handlers.NewHandler(
		jwt,
		userService,
		tagService,
		categoryService,
		articleService,
	)

	server := admin.NewServer(cfg, jwt, handler)

	done := make(chan struct{})

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		fmt.Println("Shutdown signal received, finishing in-flight requests...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			fmt.Printf("Shutdown error: %v\n", err)
		}

		close(done)
	}()

	slog.Info("Server running on", "port", cfg.AdminServer.Port)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("Server error: %v\n", err)
	}

	<-done
	fmt.Println("Server exited cleanly")
}

func loadConfig() *config.Config {
	envPath := filepath.Join(".env")

	cfg, err := config.New(envPath)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to load config %v", err))

		os.Exit(1)
	}

	return cfg
}

func loadUserService(repository *db.Repository) *service.UserService {
	userValidator := validation.NewUserValidator()

	service := service.NewUserService(repository, userValidator)

	return service
}

func loadCategoryService(repository *db.Repository) *service.CategoryService {
	categoryValidator := validation.NewCategoryValidator()

	service := service.NewCategoryService(repository, categoryValidator)

	return service
}

func loadTagService(repository *db.Repository) *service.TagService {
	tagValidator := validation.NewTagValidator()

	service := service.NewTagService(repository, tagValidator)

	return service
}

func loadArticleService(repository *db.Repository) *service.ArticleService {
	service := service.NewArticleService(repository)

	return service
}
