package main

import (
	"context"
	"log"
	"os"

	"github.com/ezeportela/go-rest-ws/handlers"
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
	r.HandleFunc("/healthcheck", handlers.HealthCheckHandler(s)).Methods("GET")
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods("POST")
}
