package router

import (
	"encoding/json"
	"net/http"

	"main.go/internal/domain"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *ArticlesHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var body CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		handleRequestError(w, err)
		return
	}

	User := domain.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password,
	}

	id, err := a.svc.CreateUser(r.Context(), &User)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"id": id})
}
