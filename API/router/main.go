package router

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"main.go/internal/domain"
)

var ErrInvalidId = errors.New("id should be an integer")

type ArticlesService interface {
	CreateUser(ctx context.Context, domUser *domain.User) (int64, error)
}

type ErrorResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

type ArticlesHandler struct {
	svc ArticlesService
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

func NewArticlesRouter(svc ArticlesService) chi.Router {
	h := &ArticlesHandler{svc: svc}

	r := chi.NewRouter()
	r.Post("/users/register", h.CreateUser)

	return r
}
