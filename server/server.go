package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseUrl string
}

type Server interface {
	Config() *Config
	Start(binder func(s Server, r *mux.Router))
}

type Broker struct {
	config *Config
	router *mux.Router
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
	return &Broker{config: config, router: router}, nil
}

func (b *Broker) Config() *Config {
	return b.config
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)

	log.Println("starting server on port", b.config.Port)
	if err := http.ListenAndServe(b.config.Port, b.router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
