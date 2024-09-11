package main

import (
	"context"
	"log"
	"os"

	"github.com/ezeportela/go-rest-ws/routes"
	"github.com/ezeportela/go-rest-ws/server"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

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

	s.Start(routes.BindRoutes)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func loadEnv() {
	if !fileExists(".env") {
		return
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
