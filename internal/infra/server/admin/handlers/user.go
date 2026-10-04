package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func (h *Handler) UserProfile(w http.ResponseWriter, r *http.Request) {
	var userID int64

	_, claims, _ := jwtauth.FromContext(r.Context())

	if val, ok := claims["user_id"].(float64); ok {
		userID = int64(val)
	}

	profile, err := h.usvc.GetProfile(r.Context(), userID)
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
