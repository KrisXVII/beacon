package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.ResponseWriter.WriteHeader(code)
	rec.status = code
}

func Log(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r) // run handler

			logger.Info(fmt.Sprintf("%s %s %d", r.Method, r.URL.Path, rec.status))
		})
	}
}
