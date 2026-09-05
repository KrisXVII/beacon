package api

import (
	"log/slog"
	"net/http"
)

// Server holds the dependencies shared by every HTTP handler.
type Server struct {
	logger *slog.Logger
}

// NewServer builds a Server with its dependencies.
func NewServer(logger *slog.Logger) *Server {
	return &Server{logger: logger}
}

// Routes returns an http.Handler with every route registered.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.healthCheck)
	mux.HandleFunc("POST /events", s.createEvent)
	return mux
}
