package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/tsaqiffatih/auth-service/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "fail",
				"error": map[string]interface{}{
					"code":    http.StatusUnauthorized,
					"message": "No token provided",
				},
			})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Verify token
		claims, err := utils.VerifyJWT(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "fail",
				"error": map[string]interface{}{
					"code":    http.StatusUnauthorized,
					"message": "Invalid token",
				},
			})
			return
		}

		// Token valid, continue with the request
		r.Header.Set("User", claims.Username)
		next.ServeHTTP(w, r)
	})
}
