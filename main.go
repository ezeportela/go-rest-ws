package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/ezeportela/go-rest-ws/handlers"
	"github.com/ezeportela/go-rest-ws/middlewares"
	"github.com/ezeportela/go-rest-ws/server"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	jwtSecret := os.Getenv("JWT_SECRET")
	databaseUrl := os.Getenv("DATABASE_URL")

	s, err := server.NewServer(context.Background(), &server.Config{
		Port:        port,
		JWTSecret:   jwtSecret,
		DatabaseUrl: databaseUrl,
	})

	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	s.Start(bindRoutes)
}

func bindRoutes(s server.Server, r *mux.Router) {
	r.Use(middlewares.AuthMiddleware(s, []string{"healthcheck", "login", "signup"}))

	r.HandleFunc("/healthcheck", handlers.HealthCheckHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)

	r.HandleFunc("/posts", handlers.CreatePostHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/posts", handlers.ListPostsHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/posts/{id}", handlers.GetPostByIdHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/posts/{id}", handlers.UpdatePostHandler(s)).Methods(http.MethodPut)
	r.HandleFunc("/posts/{id}", handlers.DeletePostHandler(s)).Methods(http.MethodDelete)
}
