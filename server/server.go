package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/ezeportela/go-rest-ws/database"
	"github.com/ezeportela/go-rest-ws/repositories"
	"github.com/ezeportela/go-rest-ws/websocket"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseUrl string
}

type Server interface {
	Config() *Config
	Hub() *websocket.Hub
	Start(binder func(s Server, r *mux.Router))
}

type Broker struct {
	config *Config
	router *mux.Router
	hub    *websocket.Hub
}

func NewServer(ctx context.Context, config *Config) (Server, error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("jwt secret is required")
	}
	if config.DatabaseUrl == "" {
		return nil, errors.New("database url is required")
	}

	router := mux.NewRouter()
	return &Broker{
		config: config,
		router: router,
		hub:    websocket.NewHub(),
	}, nil
}

func (b *Broker) Config() *Config {
	return b.config
}

func (b *Broker) Hub() *websocket.Hub {
	return b.hub
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)

	handler := cors.Default().Handler(b.router)

	repo, err := database.NewPostgresRepository(b.config.DatabaseUrl)
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}
	repositories.SetUserRepository(repo)

	go b.hub.Run()

	log.Println("starting server on port", b.config.Port)
	if err := http.ListenAndServe(b.config.Port, handler); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
