package router

import (
	"encoding/json"
	"net/http"
)

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *UserHandler) login(w http.ResponseWriter, r *http.Request) {
	var body UserLogin

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&body); err != nil {
		handleRequestError(w, err)
		return
	}

	token, err := a.svc.Login(r.Context(), body.Email, body.Password, a.jwt)
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
