package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// router set app endpoints and gets http.Handler handler.
func (s *Service) router() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Set endpoint for checking service is still alive.
	mux.Use(middleware.Heartbeat("/ping"))

	// Set other endpoints.
	mux.Post("/", s.Broker)
	mux.Post("/handle", s.HandleSubmission)
	mux.Post("/log-grpc", s.LogViaGRPC)

	return mux
}
