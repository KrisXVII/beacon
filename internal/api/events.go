package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

	err := s.notifier.Notify(event, 3)
	if err != nil {
		s.logger.Error("notify failed", "err", err.Error())
		panic(errors.New(fmt.Sprintf("failed to reach Slack: %s", err.Error())))
	}

	s.responder.JSON(w, http.StatusOK, event)
}

func (s *Server) echoEvent(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var pretty bytes.Buffer
	json.Indent(&pretty, body, "", "  ")
	fmt.Println(pretty.String())
	w.WriteHeader(http.StatusOK)
}
