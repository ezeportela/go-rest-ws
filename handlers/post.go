package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ezeportela/go-rest-ws/models"
	"github.com/ezeportela/go-rest-ws/repositories"
	"github.com/ezeportela/go-rest-ws/server"
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
		claims, err := getToken(w, r, s.Config().JWTSecret)
		if err != nil {
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
		claims, err := getToken(w, r, s.Config().JWTSecret)
		if err != nil {
			return
		}

		post, err := repositories.GetPostById(r.Context(), getPostId(r))
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
		claims, err := getToken(w, r, s.Config().JWTSecret)
		if err != nil {
			return
		}

		var req UpsertPostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		post := &models.Post{
			Id:          getPostId(r),
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
		claims, err := getToken(w, r, s.Config().JWTSecret)
		if err != nil {
			return
		}

		if err := repositories.DeletePost(r.Context(), getPostId(r), claims.UserId); err != nil {
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

func getPostId(r *http.Request) string {
	return mux.Vars(r)["id"]
}

func ListPostsHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.URL.Query().Get("limit")
		pageStr := r.URL.Query().Get("page")

		limit := uint64(10)
		if limitStr != "" {
			var err error
			limit, err = strconv.ParseUint(limitStr, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		page := uint64(0)
		if pageStr != "" {
			var err error
			page, err = strconv.ParseUint(pageStr, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		posts, err := repositories.ListPosts(r.Context(), limit, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(posts)
	}
}
