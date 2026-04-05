# go-echo-blueprint

Production-ready Echo/Go microservice template. Uses [`go-utils`](https://github.com/juanmaAV/go-utils) for observability, validation, error handling, and identity propagation.

## Stack

- **Echo v4** — HTTP framework
- **go-utils** — logger, telemetry, validator, errors, identity middleware
- **OpenTelemetry** — distributed tracing (OTLP gRPC export)
- **golangci-lint** — static analysis in CI

## Getting started

```bash
# 1. Clone and rename
cp -r go-echo-blueprint my-service
cd my-service

# 2. Update module name
find . -type f -name "*.go" -exec sed -i '' 's|go-echo-blueprint|my-service|g' {} +
# Update go.mod module line and config/config.go ServiceName

# 3. Copy env
cp .env.example .env

# 4. Run
make run
```

## Structure

```
├── config/           — configuration loaded from env vars
├── server/
│   ├── start.go      — bootstrap: config → telemetry → wire deps → run
│   ├── server.go     — Echo setup + graceful shutdown
│   └── router.go     — middleware + route registration
├── internal/
│   └── health/       — one directory per domain
│       ├── handler.go  — HTTP layer + RegisterRoutes
│       └── service.go  — business logic
└── tests/
    ├── health_test.go
    └── helpers/server.go
```

## Adding a new domain

1. Create `internal/<domain>/service.go` — define `Service` interface and implementation
2. Create `internal/<domain>/handler.go` — define `Handler` interface, `RegisterRoutes`, and implementation
3. Register in `server/router.go`:
   ```go
   // in Handlers struct
   Orders orders.Handler

   // in RegisterRoutes
   orders.RegisterRoutes(protected, h.Orders)
   ```
4. Wire in `server/start.go`:
   ```go
   handlers := Handlers{
       Health: health.NewHandler(health.NewService()),
       Orders: orders.NewHandler(orders.NewService(db, log)),
   }
   ```

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP listen port |
| `GRACEFUL_TIME` | `10s` | Graceful shutdown timeout |
| `ENVIRONMENT` | — | `local`, `staging`, `production` |
| `OTEL_EXPORTER_ENDPOINT` | — | OTLP gRPC collector (empty = disabled) |
| `OTEL_INSECURE` | `true` | Disable TLS for OTLP connection |

## Commands

```bash
make run    # go run main.go
make build  # go build -o bin/app
make test   # go test ./... -race -coverprofile=coverage.out
make lint   # golangci-lint run
make tidy   # go mod tidy
```

## Identity middleware

Protected routes use `identity.Middleware` from go-utils. Configure project-specific headers in `server/router.go`:

```go
protected := api.Group("", identity.Middleware(identity.HeaderConfig{
    Extra: []string{"X-Tenant-ID", "X-Hierarchy-Path"},
}))
```

See [`middleware/identity` docs](https://github.com/juanmaAV/go-utils/tree/main/middleware/identity).
