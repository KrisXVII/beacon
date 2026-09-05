package api

import (
	"encoding/json"
	"net/http"

	"github.com/KrisXVII/beacon/internal/alert"
)

// HTTP handlers, endpoints in this package

func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) {
	var event alert.Event

	if err := decodeJSON(r, &event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := event.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//fmt.Println(event)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // must precede the body, or Encode sends 200
	// events.go
	if err := json.NewEncoder(w).Encode(event); err != nil {
		s.logger.Error("encoding response", "handler", "createEvent", "err", err)
	}
}
