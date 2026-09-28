package router

import (
	"encoding/json"
	"net/http"
)

type UserLogin struct {
	Email    string
	Password string
}

func (a *ArticlesHandler) login(w http.ResponseWriter, r *http.Request) {
	var body UserLogin
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		handleRequestError(w, err)
		return
	}

	token, err := a.svc.Login(r.Context(), body.Email, body.Password, a.jwt)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
	})

}
