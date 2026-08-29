# Beacon

A lightweight Go microservice that watches for critical exceptions in a Flask
backend and forwards real-time alerts to Slack.

Part of a personal system that also includes a Flask backend and a Flutter
frontend. Beacon is the alerting piece: when something goes seriously wrong
in the backend, someone should find out in Slack fast, without that
notification logic living inside — and slowing down or risking — the Flask
request path itself.

## Status

🚧 Early development. Architecture and interfaces below are still being
decided — this README will get more precise as those decisions land.

## How it works

<!-- TODO: fill in once decided -->

- **Ingestion** — how does Beacon learn an exception happened?
  (e.g. Flask POSTs to an HTTP endpoint Beacon exposes / Beacon tails a log
  file / Beacon subscribes to a queue or pub-sub channel)
- **Filtering** — what makes an exception "critical" vs. routine noise, and
  is that decision made by Flask or by Beacon?
- **Delivery** — Slack, via TBD (incoming webhook vs. bot token + Slack API)
## Tech stack

- Go [version TBD]
- [HTTP framework or stdlib `net/http` — TBD]
- [Slack webhook client or Slack API library — TBD]
## Getting started

### Prerequisites

- Go [version] installed
- A Slack incoming webhook URL (or bot token, depending on delivery choice)
### Installation

```bash
git clone https://github.com/<your-username>/beacon.git
cd beacon
go mod tidy
```

### Configuration

Beacon is configured via environment variables:

| Variable            | Description                  | Required |
|---------------------|-------------------------------|----------|
| `SLACK_WEBHOOK_URL`  | Slack incoming webhook URL   | TBD      |
| ...                  | ...                           | ...      |

### Running

```bash
go run .
```

## Project structure

<!-- TODO: fill in once package layout is settled -->

```
beacon/
├── main.go
└── ...
```

## Testing

```bash
go test ./...
```

## Roadmap / open questions

- [ ] Decide ingestion mechanism (webhook vs. log tail vs. queue)
- [ ] Define "critical" exception criteria
- [ ] Choose Slack delivery method (webhook vs. bot API)
- [ ] Decide deployment target (systemd service, container, bare binary
  alongside Flask)
## License

TBD
 
