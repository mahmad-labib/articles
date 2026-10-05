package admin

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mahmad-labib/articles/internal/config"
	"github.com/mahmad-labib/articles/internal/infra/server/admin/middlewares"
)

type Handler interface {
	AuthHandler
	UserHandler
	CategoryHandler
	TagHandler
	ArticleHandler
}

func NewServer(cfg *config.Config, jwt *jwtauth.JWTAuth, h Handler) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/v1/", v1routers(jwt, h))

	server := &http.Server{
		Addr:         ":" + cfg.AdminServer.Port,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.AdminServer.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.AdminServer.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.AdminServer.IdleTimeout) * time.Second,
	}

	return server
}

func v1routers(jwt *jwtauth.JWTAuth, h Handler) chi.Router {
	r := chi.NewRouter()

	// public routes
	r.Group(func(r chi.Router) {
		r.Post("/auth/register", h.Register)
		r.Post("/auth/login", h.Login)
	})

	// private routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(jwt))
		r.Use(middlewares.AuthMiddleware())

		r.Route("/admin", func(r chi.Router) {
			// user routes
			r.Get("/user/profile", h.UserProfile)

			// category routes
			r.Post("/categories", h.CreateCategory)
			r.Get("/categories", h.ListCategories)
			r.Get("/categories/{id}", h.GetCategoryByID)
			r.Put("/categories/{id}", h.UpdateCategory)
			r.Delete("/categories/{id}", h.DeleteCategory)

			// tag routes
			r.Post("/tags", h.CreateTag)
			r.Get("/tags", h.ListTags)
			r.Get("/tags/{id}", h.GetTagByID)
			r.Put("/tags/{id}", h.UpdateTag)
			r.Delete("/tags/{id}", h.DeleteTag)

			// article routes
			r.Post("/articles", h.CreateArticle)
			r.Get("/articles/{id}", h.GetArticleByID)
		})
	})

	return r
}
