package api

import (
	"log/slog"
	"net/http"

	"github.com/KrisXVII/beacon/internal/httpx"
)

// Server holds the dependencies shared by every HTTP handler.
type Server struct {
	//logger    *slog.Logger
	responder *httpx.Responder
}

// NewServer builds a Server with its dependencies. Called in main.go to create a server to configure
func NewServer(logger *slog.Logger) *Server {
	return &Server{ // &Server to return a pointer, without the & every call creates a new Server object
		//logger:    logger,
		responder: httpx.New(logger),
	}
}

// Routes returns an http.Handler with every route registered.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.healthCheck)
	mux.HandleFunc("POST /events", s.createEvent)
	return mux
}
