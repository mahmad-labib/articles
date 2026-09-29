package router

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"gorm.io/gorm"
	"main.go/internal/domain"
	"main.go/internal/service"
	"main.go/internal/service/validation"
)

type UserService interface {
	CreateUser(ctx context.Context, domUser *domain.User) (int64, error)
	Login(ctx context.Context, email, password string, jwt *jwtauth.JWTAuth) (string, error)
	GetProfile(ctx context.Context, id int64) (*domain.Profile, error)
}

type ErrorResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

type UserHandler struct {
	svc UserService
	jwt *jwtauth.JWTAuth
}

func handleServiceError(w http.ResponseWriter, err error) {
	// only for development
	fmt.Println(err)

	switch {
	case errors.Is(err, validation.EmailErr),
		errors.Is(err, validation.Empty),
		errors.Is(err, validation.NameLength),
		errors.Is(err, validation.EmailUnique),
		errors.Is(err, service.PasswordInvalid),
		errors.Is(err, gorm.ErrInvalidValue),
		errors.Is(err, validation.PasswordSpaces):
		writeError(w, http.StatusUnprocessableEntity, "Validation Error", []string{err.Error()})
	case errors.Is(err, gorm.ErrRecordNotFound):
		writeError(w, http.StatusNotFound, "Not Found", []string{err.Error()})

	default:
		writeError(w, http.StatusInternalServerError, "internal server error", []string{})
	}

}

func handleRequestError(w http.ResponseWriter, err error) {
	// only for development
	fmt.Println(err)

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

func NewUsersRouter(svc UserService, jwt *jwtauth.JWTAuth) chi.Router {
	h := &UserHandler{svc: svc, jwt: jwt}

	r := chi.NewRouter()
	r.Post("/users/register", h.CreateUser)
	r.Post("/users/login", h.login)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(jwt))
		r.Use(Authenticator())

		r.Get("/user/profile", h.UserProfile)
	})

	return r
}

func Authenticator() func(http.Handler) http.Handler {
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
