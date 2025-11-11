## CAATSM Dashboard

CAATSM Dashboard is a lightweight Go application that renders a web UI for monitoring civil-aviation telegram traffic. It combines an HTTP server (Echo), HTML templating (`github.com/a-h/templ`), and an optional NATS subscriber (Watermill) to receive and display telegram data.

- `main.go` bootstraps the Echo server, serves static assets under `/public`, and wires HTTP routes.
- `handlers` defines Echo handlers and shared helpers for rendering templ components.
- `views` holds templ templates. `Base.templ` defines the shared layout while `index.templ` renders a placeholder dashboard.
- `internal/config` loads strongly typed configuration from `configs/config.<env>.toml`.
- `internal/nats` provides a Watermill-based subscriber that streams telegram messages from NATS to user-defined handlers.
- `pkg/utils/log.go` exposes a Zap logger with environment-aware configuration.

## Prerequisites

- Go 1.25+
- Node.js/npm or pnpm (for managing static dependencies such as Flowbite and Tailwind)
- Optional: a running NATS server if you need real telegram ingestion

## Getting Started

```bash
# Install Go dependencies
go mod download

# Install frontend dependencies (used for styles and supporting assets)
npm install
```

### Environment variables

Most runtime options are driven through files in `configs/` and can be overridden with environment variables (Viper is configured with the `TELE_` prefix). Create a `.env` if you want local defaults, otherwise export the variables directly:

```
LISTEN_ADDR=:3002
GO_ENV=dev
TELE_MODE=dev
TELE_HASURA_SECRET=replace-me
```

Note: `.env` is optional—missing files emit a warning but do not abort startup.

### Configuration files

- `configs/config.dev.toml` – development defaults for server port, NATS, subscription topic, and timeout values.
- `configs/logger.dev.json` – Zap logger settings.
- `event.toml` – sample event bus configuration retained for reference.

> ⚠️ Do not commit production secrets. Override `hasura.secret` with the `TELE_HASURA_SECRET` environment variable (the default file purposely leaves it blank).

## Running the Server

```bash
go run .
```

The dashboard will be available at `http://localhost:3002/` (or whatever `LISTEN_ADDR` specifies). Static assets are served from `public/`.

## NATS Subscriber (Optional)

`internal/nats/sub.go` wraps Watermill’s NATS integration and now executes handlers via a small worker pool to avoid blocking the subscription loop. To consume messages:

1. Load configuration using `internal/config.LoadConfig()`.
2. Construct a subscriber with `nats.NewSub(config)`.
3. Implement `iface.MessageHandler` and pass it to `Subscribe`, providing a `context.Context` for lifecycle management.

The subscriber logs via `pkg/utils` and lets your handler decide how to process payloads (e.g., render them in the dashboard, persist to storage, emit metrics).

## Development Workflow

- Update templ components under `views/` and rerun the server. Templ recompilation occurs on `go run` or `go build`.
- Customize Tailwind/Flowbite assets under `public/` or via build tools referenced in `package.json`/`tailwind.config.js`.
- Consider adding tests around configuration loading and NATS integration with `go test ./...`.

## Project Structure

```
.
├── configs/              # TOML & logger settings
├── handlers/             # Echo route setup and render helpers
├── internal/
│   ├── config/           # Viper-backed configuration loader
│   ├── iface/            # Interfaces for messaging components
│   └── nats/             # Watermill-based NATS subscriber
├── pkg/utils/            # Shared utilities (logging)
├── public/               # Static assets (icons, styles, scripts)
├── views/                # templ components used by the dashboard
└── main.go               # Application entry point
```

## Troubleshooting

- **Missing `.env`** – Startup logs a warning; export the required `TELE_` variables or create a local `.env`.
- **NATS connection errors** – Ensure the URL, credentials, and TLS settings in `config.<env>.toml` (or corresponding env vars) match your NATS deployment.
- **Template rendering failures** – Echo now reports HTTP 500 with a logged error when rendering fails; check templ components and handler logic.

## Next Steps

- Flesh out the dashboard by wiring real telegram message handlers.
- Harden configuration (environment overrides, production secrets management, TLS).
- Add automated tests for handlers, config loading, and message processing (`go test ./...` provides a baseline).

