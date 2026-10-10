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
	commentService := loadCommentService(repository)

	//article publish service
	publishDone := make(chan struct{})
	go func() {
		ticker := time.NewTicker(15 * time.Minute)
		defer ticker.Stop()
		defer close(publishDone)

		publishWorker(ticker, publishDone, articleService)
	}()

	//channel handels the views count
	viewCounter := service.NewBatchCounter[int64](1000, 5*time.Second, func(counts map[int64]int) error {
		for id, views := range counts {
			if err := articleService.AddView(context.Background(), id, views); err != nil {
				return err
			}
		}
		return nil
	})

	go viewCounter.Run()

	handler := handlers.NewHandler(
		viewCounter,
		jwt,
		userService,
		tagService,
		categoryService,
		articleService,
		commentService,
	)

	server := admin.NewServer(cfg, jwt, handler)

	done := make(chan struct{})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan

		fmt.Println("Shutdown signal received, finishing in-flight requests...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			fmt.Printf("Shutdown error: %v\n", err)
		} else {
			close(publishDone)
			<-viewCounter.Close()
		}

		close(done)
	}()

	slog.Info("Server running on", "port", cfg.AdminServer.Port)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("Server error: %v\n", err)
		sigChan <- syscall.SIGTERM
	}

	<-done
	fmt.Println("Server exited cleanly")
}

var ticker = time.NewTicker(5000 * time.Millisecond)

func publishWorker(ticker *time.Ticker, done chan struct{}, articleService *service.ArticleService) {
	for {
		select {
		case <-ticker.C:
			err := articleService.PublishArticle(context.Background())
			if err != nil {
				slog.Error("Error publishing articles:", "error", err)
			}
		case <-done:
			err := articleService.PublishArticle(context.Background())
			if err != nil {
				slog.Error("Error publishing articles:", "error", err)
			}
			return
		}
	}
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
	articleValidator := validation.NewArticleValidator()

	service := service.NewArticleService(repository, articleValidator)

	return service
}

func loadCommentService(repository *db.Repository) *service.CommentService {
	commentValidator := validation.NewCommentValidator()

	service := service.NewCommentService(repository, commentValidator)

	return service
}
