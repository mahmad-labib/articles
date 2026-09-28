package router

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"main.go/internal/domain"
)

var ErrInvalidId = errors.New("id should be an integer")

type ArticlesService interface {
	CreateUser(ctx context.Context, domUser *domain.User) (int64, error)
	Login(ctx context.Context, email, password string, jwt *jwtauth.JWTAuth) (string, error)
	GetProfile(ctx context.Context, id int64) (*domain.Profile, error)
}

type ErrorResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

type ArticlesHandler struct {
	svc ArticlesService
	jwt *jwtauth.JWTAuth
}

func handleServiceError(w http.ResponseWriter, err error) {
	writeError(w, http.StatusUnprocessableEntity, "Validation Error", []string{err.Error()})
}

func writeError(w http.ResponseWriter, status int, message string, details []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(ErrorResponse{
		Code:    status,
		Message: message,
		Details: details,
	})
}

func handleRequestError(w http.ResponseWriter, err error) {
	writeError(w, http.StatusBadRequest, "Bad Request", []string{err.Error()})
}

func NewArticlesRouter(svc ArticlesService, jwt *jwtauth.JWTAuth) chi.Router {
	h := &ArticlesHandler{svc: svc, jwt: jwt}

	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(jwt))
	r.Post("/users/register", h.CreateUser)
	r.Post("/users/login", h.login)
	r.Get("/user/profile", h.UserProfile)

	return r
}
