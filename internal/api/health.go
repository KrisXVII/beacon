package api

import (
	"encoding/json"
	"net/http"
)

// Health check and other top-level utils

// healthCheck reports that the service is running.
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) { // w is where the response is written, r the incoming request
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		s.logger.Error("encoding response", "handler", "healthCheck", "err", err)
	}
}
