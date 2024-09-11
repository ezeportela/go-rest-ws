package routes

import (
	"net/http"

	"github.com/ezeportela/go-rest-ws/handlers"
	"github.com/ezeportela/go-rest-ws/middlewares"
	"github.com/ezeportela/go-rest-ws/server"
	"github.com/gorilla/mux"
)

func BindRoutes(s server.Server, r *mux.Router) {
	api := r.PathPrefix("/api/v1").Subrouter()

	api.Use(middlewares.AuthMiddleware(s, []string{"healthcheck", "login", "signup"}))

	r.HandleFunc("/healthcheck", handlers.HealthCheckHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	api.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)

	api.HandleFunc("/posts", handlers.CreatePostHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/posts", handlers.ListPostsHandler(s)).Methods(http.MethodGet)
	api.HandleFunc("/posts/{id}", handlers.GetPostByIdHandler(s)).Methods(http.MethodGet)
	api.HandleFunc("/posts/{id}", handlers.UpdatePostHandler(s)).Methods(http.MethodPut)
	api.HandleFunc("/posts/{id}", handlers.DeletePostHandler(s)).Methods(http.MethodDelete)

	r.HandleFunc("/ws", s.Hub().HandleWebsocket)
}
