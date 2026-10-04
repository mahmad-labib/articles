package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func AuthMiddleware() func(http.Handler) http.Handler {
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
