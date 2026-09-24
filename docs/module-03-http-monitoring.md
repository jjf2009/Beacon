# Module 3 — HTTP Monitoring

## Status: Done

## What You're Building

Beacon sends HTTP requests to registered endpoints and records results. First time Beacon actually monitors something.

## Concepts to Learn First

- **Go HTTP client** — `http.Client`, `http.Get()`
- **HTTP status codes** — 200 OK, 5xx server error, etc.
- **Request timeouts** — don't wait forever for a response
- **Latency measurement** — `time.Now()` before/after request
- **DNS failures** — what happens when hostname doesn't resolve
- **Connection failures** — host unreachable, connection refused
- **Context cancellation** — cancel in-flight requests cleanly

## What a Check Records

```
status        UP | DOWN
status_code   200, 500, 0 (connection failed)
response_time 142 (milliseconds)
error         "" | "connection refused" | "timeout"
checked_at    2026-01-01T14:00:00Z
```

## What to Build

Single function: `CheckEndpoint(url string) CheckResult`

```go
type CheckResult struct {
    Status       string
    StatusCode   int
    ResponseTime int  // ms
    Error        string
    CheckedAt    time.Time
}
```

Logic:
1. Record start time
2. Send GET request with timeout (e.g. 10s)
3. Record end time → compute latency
4. If error → status=DOWN, store error message
5. If status code >= 400 → status=DOWN
6. Else → status=UP
7. Save CheckResult to `checks` table

## Timeout Config

```go
client := &http.Client{
    Timeout: 10 * time.Second,
}
```

Never use `http.Get()` directly — it has no timeout. Always use custom client.

## Error Cases to Handle

| Scenario | What to record |
|----------|---------------|
| 200 response | UP, 200, latency |
| 500 response | DOWN, 500, latency |
| Connection refused | DOWN, 0, "connection refused" |
| DNS failure | DOWN, 0, "no such host" |
| Timeout | DOWN, 0, "timeout" |

## File Structure Target

```
backend/
  internal/
    checker/
      checker.go     ← CheckEndpoint function
    check/
      model.go       ← Check struct
      repository.go  ← save check to DB
```

## Add API Endpoint to Trigger a Check

```
POST /api/endpoints/:id/check   ← manually trigger one check
GET  /api/endpoints/:id/checks  ← list recent checks
```

This lets you test checking without the background worker (Module 4).

## Test Commands

```bash
# Trigger check on endpoint
curl -X POST http://localhost:8080/api/endpoints/<id>/check

# List checks for endpoint
curl http://localhost:8080/api/endpoints/<id>/checks

# Test against a DOWN endpoint
# Stop a local server and check — should record DOWN
```

## Definition of Done

- [ ] `CheckEndpoint` returns correct result for UP endpoint
- [ ] `CheckEndpoint` returns DOWN on connection refused
- [ ] `CheckEndpoint` returns DOWN on timeout
- [ ] Results saved to `checks` table
- [ ] `GET /api/endpoints/:id/checks` returns check history
- [ ] Timeout is configurable, not hardcoded
- [ ] Can explain what happens at network level during a check

## Questions Before Moving On

1. Why set a timeout on HTTP client?
2. What is the difference between a DNS failure and connection refused?
3. How do you measure latency accurately?
4. What does HTTP status code 0 mean in your system?

## Next: Module 4 — Background Worker
