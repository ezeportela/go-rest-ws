package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ezeportela/go-rest-ws/models"
	"github.com/ezeportela/go-rest-ws/repositories"
	"github.com/ezeportela/go-rest-ws/server"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

type UpsertPostRequest struct {
	PostContent string `json:"post_content"`
}

type InsertPostResponse struct {
	Id          string `json:"id"`
	PostContent string `json:"post_content"`
}

type ActionResponse struct {
	Message string `json:"message"`
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

		var req UpsertPostRequest
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

func GetPostByIdHandler(s server.Server) http.HandlerFunc {
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

		params := mux.Vars(r)
		post, err := repositories.GetPostById(r.Context(), params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if post == nil {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}

		if post.UserId != claims.UserId {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(post)
	}
}

func UpdatePostHandler(s server.Server) http.HandlerFunc {
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

		var req UpsertPostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		params := mux.Vars(r)
		post := &models.Post{
			Id:          params["id"],
			UserId:      claims.UserId,
			PostContent: req.PostContent,
		}

		if err := repositories.UpdatePost(r.Context(), post); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(&ActionResponse{
			Message: "post updated",
		})
	}
}

func DeletePostHandler(s server.Server) http.HandlerFunc {
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

		params := mux.Vars(r)
		if err := repositories.DeletePost(r.Context(), params["id"], claims.UserId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(&ActionResponse{
			Message: "post deleted",
		})
	}
}
