package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/KrisXVII/beacon/internal/api"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	server := api.NewServer(logger)

	logger.Info("beacon listening", "addr", ":8080")
	if err := http.ListenAndServe(":8080", server.Routes()); err != nil {
		logger.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
