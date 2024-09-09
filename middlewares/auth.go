package middlewares

import (
	"net/http"
	"strings"

	"github.com/ezeportela/go-rest-ws/models"
	"github.com/ezeportela/go-rest-ws/server"
	"github.com/golang-jwt/jwt"
)

func shouldCheckAuth(route string, noAuthNeeded []string) bool {
	for _, path := range noAuthNeeded {
		if strings.Contains(route, path) {
			return false
		}
	}
	return true
}

func AuthMiddleware(s server.Server, noAuthNeeded []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !shouldCheckAuth(r.URL.Path, noAuthNeeded) {
				next.ServeHTTP(w, r)
				return
			}

			tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
			_, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JWTSecret), nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
