package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mahmad-labib/articles/internal/domain"
)

type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
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

	id, err := h.usvc.CreateUser(r.Context(), &user)
	if err != nil {
		handleServiceError(w, err)

		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"id": id})
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var body UserLogin

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		handleRequestError(w, err)

		return
	}

	token, err := h.usvc.Login(r.Context(), body.Email, body.Password, h.jwt)
	if err != nil {
		handleServiceError(w, err)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
	})
}
