package middlewares

import (
	"encoding/json"
	"net/http"
	"strings"
)

var (
	NO_AUTH_NEEDED = []string{
		"login",
		"signup",
	}
)

func shoulCheckToken(route string) bool {
	for _, p := range NO_AUTH_NEEDED {
		if strings.Contains(route, p) {
			return false
		}
	}
	return true
}

func CheckAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !shoulCheckToken(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		userID := r.Header.Get("Authorization")
		if userID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "you are not logged in"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
