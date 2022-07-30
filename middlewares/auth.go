package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/DarioRoman01/cqrs/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
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

func CheckAuthMiddleware(secret string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !shoulCheckToken(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			claims := &models.Claims{}
			tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
			_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", claims.UserID)))
		})
	}
}
