package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ezeportela/go-rest-ws/models"
	"github.com/ezeportela/go-rest-ws/repositories"
	"github.com/ezeportela/go-rest-ws/server"
	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
)

type InsertPostRequest struct {
	PostContent string `json:"post_content"`
}

type InsertPostResponse struct {
	Id          string `json:"id"`
	PostContent string `json:"post_content"`
}

func CreatePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*models.AppClaims)
		if !ok || !token.Valid {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		var req InsertPostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		post := &models.Post{
			Id:          id.String(),
			UserId:      claims.UserId,
			PostContent: req.PostContent,
		}

		if err := repositories.InsertPost(r.Context(), post); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(&InsertPostResponse{
			Id:          post.Id,
			PostContent: post.PostContent,
		})
	}
}
