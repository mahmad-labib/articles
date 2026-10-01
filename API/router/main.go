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
	CreateCategory(ctx context.Context, category *domain.Category) (int64, error)
	ListCategories(ctx context.Context, page, pageSize int) ([]domain.Category, error)
	GetCategoryByID(ctx context.Context, id int64) (*domain.Category, error)
	UpdateCategory(ctx context.Context, category *domain.Category) error
	DeleteCategory(ctx context.Context, id int64) error
	CreateTag(ctx context.Context, tag *domain.Tag) (int64, error)
	ListTags(ctx context.Context, page, pageSize int) ([]domain.Tag, error)
	GetTagByID(ctx context.Context, id int64) (*domain.Tag, error)
	UpdateTag(ctx context.Context, tag *domain.Tag) error
	DeleteTag(ctx context.Context, id int64) error
}

var ErrInvalidId = errors.New("id should be an integer")

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
	// ----------only for development------------
	fmt.Println(err)
	fmt.Println(errors.Is(err, gorm.ErrDuplicatedKey))

	switch {
	case errors.Is(err, validation.EmailErr),
		errors.Is(err, validation.Empty),
		errors.Is(err, validation.NameLength),
		errors.Is(err, validation.EmailUnique),
		errors.Is(err, service.PasswordInvalid),
		errors.Is(err, gorm.ErrInvalidValue),
		errors.Is(err, validation.FieldSpaces):
		writeError(w, http.StatusUnprocessableEntity, "Validation Error", []string{err.Error()})

	case errors.Is(err, gorm.ErrRecordNotFound):
		writeError(w, http.StatusNotFound, "Not Found", []string{err.Error()})

	case errors.Is(err, gorm.ErrDuplicatedKey):
		writeError(w, http.StatusConflict, "item already exist", []string{err.Error()})

	default:
		writeError(w, http.StatusInternalServerError, "internal server error", []string{})
	}

}

func handleRequestError(w http.ResponseWriter, err error) {
	// ---------only for development------------
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
		r.Post("/categories", h.CreateCategory)
		r.Get("/categories", h.ListCategories)
		r.Get("/categories/{id}", h.GetCategoryByID)
		r.Put("/categories/{id}", h.UpdateCategory)
		r.Delete("/categories/{id}", h.DeleteCategory)
		r.Post("/tags", h.CreateTag)
		r.Get("/tags", h.ListTags)
		r.Get("/tags/{id}", h.GetTagByID)
		r.Put("/tags/{id}", h.UpdateTag)
		r.Delete("/tags/{id}", h.DeleteTag)
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
