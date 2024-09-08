package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ezeportela/go-rest-ws/server"
)

type HealthCheckResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func HealthCheckHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(HealthCheckResponse{
			Message: "Server is running",
			Status:  true,
		})
	}
}
