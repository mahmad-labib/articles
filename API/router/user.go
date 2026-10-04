package router

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"main.go/internal/domain"
)

type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var body UserRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		handleRequestError(w, err)

		return
	}

	user := domain.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password,
	}

	id, err := a.svc.CreateUser(r.Context(), &user)
	if err != nil {
		handleServiceError(w, err)

		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"id": id})
}

func (a *UserHandler) UserProfile(w http.ResponseWriter, r *http.Request) {
	var userID int64

	_, claims, _ := jwtauth.FromContext(r.Context())

	if val, ok := claims["user_id"].(float64); ok {
		userID = int64(val)
	}

	profile, err := a.svc.GetProfile(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"profile": profile,
	})
}
