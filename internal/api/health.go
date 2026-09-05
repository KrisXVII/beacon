package api

import (
	"net/http"
)

// Health check and other top-level utils

// healthCheck reports that the service is running.
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) { // w is where the response is written, r the incoming request
	s.responder.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
