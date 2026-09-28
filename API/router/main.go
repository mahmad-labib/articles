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

func handleRequestError(w http.ResponseWriter, err error) {
	writeError(w, http.StatusBadRequest, "Bad Request", []string{err.Error()})
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

func NewArticlesRouter(svc ArticlesService, jwt *jwtauth.JWTAuth) chi.Router {
	h := &ArticlesHandler{svc: svc, jwt: jwt}

	r := chi.NewRouter()
	r.Post("/users/register", h.CreateUser)
	r.Post("/users/login", h.login)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(jwt))
		r.Use(Authenticator(jwt))

		r.Get("/user/profile", h.UserProfile)
	})

	return r
}

func Authenticator(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())

			w.Header().Set("Content-Type", "application/json")

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}

			if token == nil {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized: token not found"})
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}
