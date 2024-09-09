package handlers

import (
	"net/http"
	"strings"

	"github.com/ezeportela/go-rest-ws/models"
	"github.com/golang-jwt/jwt"
)

func getToken(w http.ResponseWriter, r *http.Request, secret string) (*models.AppClaims, error) {
	tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
	token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return nil, err
	}

	claims, ok := token.Claims.(*models.AppClaims)
	if !ok || !token.Valid {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return nil, err
	}

	return claims, nil
}
