# Module 15 — Production Hardening

## Status: NOT STARTED

## What You're Building

Make Beacon deployable as a real service. Not just "works on my machine." Proper containers, config management, CI/CD, and reliability.

## Topics to Cover

### 1. Docker

Containerize every service:

```dockerfile
# backend/Dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o api ./cmd/api

FROM alpine:latest
COPY --from=builder /app/api /api
EXPOSE 8080
CMD ["/api"]
```

```dockerfile
# frontend/Dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:20-alpine
COPY --from=builder /app/.next ./.next
CMD ["npm", "run", "start"]
```

Multi-stage build: builder stage compiles, final stage is small.

### 2. Health Checks

Every container needs a health check so orchestrators know when it's ready:

```yaml
# docker-compose.yml
api:
  healthcheck:
    test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
    interval: 10s
    timeout: 5s
    retries: 3
```

### 3. Graceful Shutdown

Handle SIGTERM — finish in-flight requests before exiting. You already did most of this in your `main.go`. As an exercise, reconstruct it from memory:

1. Run `server.ListenAndServe()` in a goroutine (it blocks — you don't want it blocking main). Treat `http.ErrServerClosed` as normal, not an error.
2. Block on `<-ctx.Done()` (the signal context you built with `signal.NotifyContext`).
3. Once cancelled, make a new context with a timeout (say 30s) so shutdown can't hang forever.
4. Call `server.Shutdown(shutdownCtx)` — it stops accepting new requests and waits for in-flight ones to finish.

Question: what happens to a request that's mid-flight if you skip step 4 and just `os.Exit`?

### 4. Structured Logging

Replace `fmt.Println` with structured logs (JSON) using `log/slog`. You've been doing this already.

**Exercise:** write a `slog.Info` call that logs "check completed" with three key/value fields: the endpoint ID, the status, and the latency. Recall the pattern: `slog.Info(message, key1, val1, key2, val2, ...)` — pairs of key then value.

Output:
```json
{"time":"2026-01-01T14:00:00Z","level":"INFO","msg":"check completed","endpoint_id":"abc123","status":"up","latency_ms":142}
```

Machine-parseable. Log aggregation systems (Loki, Datadog) can query by field.

### 5. Configuration Management

All config from environment variables. No hardcoded values. No config files checked into git.

**Exercise:** you already have a `Config` struct and loader. Extend it for production. Add fields for anything still hardcoded (Redis URL, API key, etc.). Write two small helpers — `getEnv(key, fallback string) string` and `getEnvInt(key string, fallback int) int` — that read an env var and return a default when it's unset. In `Load`, build the `Config` from those helpers. Think: why default values? What should happen if a *required* secret like the DB password is missing — default, or crash loudly?

### 6. Secrets Management

Never commit secrets. Use:
- `.env` file locally (in `.gitignore`)
- Environment variables in production
- GitHub Actions secrets for CI

```bash
# .env (never committed)
DB_PASSWORD=supersecret
DISCORD_WEBHOOK_URL=https://discord.com/...
BEACON_API_KEY=random-string-here
```

### 7. CI/CD with GitHub Actions

```yaml
# .github/workflows/ci.yml
name: CI

on: [push, pull_request]

jobs:
  backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.23' }
      - run: go build ./...
      - run: go vet ./...
      - run: go test ./...

  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: '20' }
      - run: npm ci
      - run: npm run build
      - run: npm run lint
```

### 8. Database Migrations (automated)

Stop running migrations by hand. Use a migration tool:

```bash
# Install golang-migrate
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run on startup or via CI
migrate -path ./migrations -database $DATABASE_URL up
```

### 9. Rate Limiting

Protect API from abuse.

**Exercise:** use `golang.org/x/time/rate`. Create a limiter that allows ~10 requests/second (`rate.NewLimiter`). Then write `rateLimitMiddleware` — same middleware shape you learned in Module 10 (takes an `http.Handler`, returns one). Inside the wrapped handler: if `limiter.Allow()` returns false, respond with 429 (`http.StatusTooManyRequests`) and `return`; otherwise call `next.ServeHTTP`. This is the middleware pattern again — once you know it, every cross-cutting concern (auth, logging, rate limiting) follows the same shape.

### 10. Integration Tests

Test the full stack: API + real database. Not unit tests with mocks.

**Exercise:** write `TestCreateEndpoint(t *testing.T)`. A Go test is a function starting with `Test`, taking `*testing.T`. The steps: start a test DB (testcontainers, or a pre-seeded local one), run migrations, call `POST /api/endpoints`, then assert both that the response is correct and that a row actually landed in the DB. Use `t.Fatalf` when a check fails. Learning to write a real test is a hiring signal on its own — many juniors can't.

## Production Deployment Options

Pick one to learn with:

| Option | Complexity | Cost | Good for |
|--------|-----------|------|----------|
| Fly.io | Low | Free tier | Learning |
| Railway | Low | Free tier | Learning |
| DigitalOcean | Medium | ~$12/mo | Real deploy |
| AWS ECS | High | Variable | Enterprise pattern |

**Recommendation:** Fly.io or Railway to learn deployment without infrastructure overhead.

## Checklist: Production Ready?

- [ ] All services run in Docker containers
- [ ] `docker compose up` starts full stack
- [ ] No secrets in git history
- [ ] Graceful shutdown works (Ctrl+C exits cleanly)
- [ ] Structured JSON logging
- [ ] Health check endpoint returns 200 only when ready
- [ ] CI passes on every push
- [ ] DB migrations run automatically
- [ ] Rate limiting on public endpoints
- [ ] Integration tests exist

## Questions Before Calling It Done

1. What is multi-stage Docker build and why is it smaller?
2. What is graceful shutdown and what happens without it?
3. What is the difference between structured and unstructured logs?
4. What is a database migration and why automate it?
5. What does CI catch that local development doesn't?

## Final Goal

After this module, Beacon is:

```
Monitored    → checks your apps
Alerted      → notifies on incidents
Correlated   → links incidents to deploys
Observable   → Prometheus + Grafana + traces
Deployable   → Docker, CI/CD, production config
```

That is the difference between a portfolio project and an engineering project.
