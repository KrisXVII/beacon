package api

import (
	"net/http"

	"github.com/KrisXVII/beacon/internal/alert"
	"github.com/KrisXVII/beacon/internal/httpx"
)

// HTTP handlers, endpoints in this package

func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) {
	var event alert.Event

	if err := httpx.Decode(r, &event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := event.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.responder.JSON(w, http.StatusOK, event)
}
