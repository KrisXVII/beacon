package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Responder struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *Responder {
	return &Responder{logger: logger}
}

func (r *Responder) JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		r.logger.Error("encoding response", "err", err)
	}
}
