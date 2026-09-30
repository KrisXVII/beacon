package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/KrisXVII/beacon/internal/alert"
)

// Link Slack API and write method to send data to Slack to create a notification on the dedicated channel

type Client struct {
	webhookURL string
	http       *http.Client
}

func New(webhookURL string) *Client {
	return &Client{
		webhookURL: webhookURL,
		http:       &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *Client) Notify(event alert.Event, count int) error {
	message, err := json.Marshal(buildMessage(event, count))
	if err != nil {
		return fmt.Errorf("encoding slack message: %w", err)
	}

	resp, slackErr := c.http.Post(c.webhookURL, "application/json", bytes.NewReader(message))
	if slackErr != nil {
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			slackErr = urlErr.Err
		}
		return fmt.Errorf("posting to slack: %w", slackErr)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("slack returned %d: %s", resp.StatusCode, body)
	}

	return nil
}

var mrkdwnEscaper = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")

func buildMessage(e alert.Event, count int) map[string]any {
	message := e.Message
	if message == "" {
		message = "(no message)"
	}

	//where := fmt.Sprintf("`%s %s`", e.Method, e.Route)
	//if f, ok := e.AppFrame(); ok {
	//	where += fmt.Sprintf(" · `%s:%d in %s`", f.File, f.Line, f.Func)
	//}

	blocks := []any{
		map[string]any{
			"type": "header",
			"text": map[string]any{"type": "plain_text", "text": "🔴 " + e.ErrorType},
		},
		map[string]any{
			"type": "section",
			"text": map[string]any{
				"type": "mrkdwn",
				"text": "*" + mrkdwnEscaper.Replace(message) + "*\n", // + mrkdwnEscaper.Replace(where),
			},
		},
		map[string]any{
			"type": "section",
			"fields": []any{
				map[string]any{"type": "mrkdwn", "text": "*Environment*\n" + e.Environment},
				map[string]any{"type": "mrkdwn", "text": fmt.Sprintf("*Occurrences*\n%d", count)},
			},
		},
		map[string]any{
			"type": "context",
			"elements": []any{
				map[string]any{
					"type": "mrkdwn",
					"text": e.Service + " · " + e.OccurredAt.UTC().Format("2006-01-02 15:04 UTC"),
				},
			},
		},
	}

	if e.PosthogEventID != "" {
		blocks = append(blocks, map[string]any{
			"type": "actions",
			"elements": []any{
				map[string]any{
					"type": "button",
					"text": map[string]any{"type": "plain_text", "text": "View in PostHog"},
					//"url":  postHogErrorsURL,
				},
			},
		})
	}

	return map[string]any{
		"text":   e.ErrorType + ": " + message, // plain-text fallback for notifications
		"blocks": blocks,
	}
}
