package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/KrisXVII/beacon/internal/adapters/slack"
	"github.com/KrisXVII/beacon/internal/api"
)

func main() {
	slackUrl := os.Getenv("SLACK_WEBHOOK_URL")
	notifier := slack.New(slackUrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	server := api.NewServer(logger, notifier)

	logger.Info("beacon listening", "addr", ":8080")
	if err := http.ListenAndServe(":8080", server.Routes()); err != nil {
		logger.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
