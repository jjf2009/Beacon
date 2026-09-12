# Module 4 — Background Worker

## Status: NOT STARTED

## What You're Building

Separate worker process that continuously checks endpoints. The API handles user requests. The worker does monitoring. Two jobs, two processes.

## Why Separate?

If the API handled monitoring, every HTTP request would block while checking endpoints. Bad. Worker runs independently — API stays fast.

```
Go API     → handles user requests (CRUD)
Go Worker  → runs health checks in background
```

## Concepts to Learn First

- **Goroutines** — lightweight threads in Go (`go func()`)
- **Channels** — communicate between goroutines
- **Concurrency** — multiple things happening at once
- **Background processes** — long-running loops
- **Graceful shutdown** — stop cleanly when OS sends SIGTERM
- **Worker patterns** — fetch work → do work → repeat

## Worker Loop (simple version)

```
1. Fetch all endpoints from DB
2. For each endpoint: run CheckEndpoint()
3. Save result to DB
4. Wait (sleep)
5. Repeat
```

## What to Build

### Option A: Worker as goroutine inside same binary

```go
// main.go
go worker.Start(db)
server.Start()
```

### Option B: Separate binary (recommended for learning)

```
backend/
  cmd/api/main.go       ← API server
  cmd/worker/main.go    ← background worker
```

Separate binary is cleaner. Forces you to understand both processes.

## Graceful Shutdown

Worker must stop cleanly on SIGTERM/SIGINT (Ctrl+C).

```go
ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
defer cancel()

// worker loop checks ctx.Done()
for {
    select {
    case <-ctx.Done():
        // cleanup, exit
        return
    default:
        // do check work
    }
}
```

## File Structure Target

```
backend/
  cmd/
    api/main.go        ← API server (existing)
    worker/main.go     ← worker entrypoint
  internal/
    worker/
      worker.go        ← worker loop logic
```

## Concurrency (stretch goal for this module)

Check endpoints concurrently instead of sequentially:

```go
for _, endpoint := range endpoints {
    go func(e Endpoint) {
        result := checker.Check(e.URL)
        repo.SaveCheck(e.ID, result)
    }(endpoint)
}
```

Be careful — if you have 1000 endpoints, spawning 1000 goroutines at once can overwhelm resources. Use a worker pool with buffered channels to limit concurrency. (Introduce this concept here but implement simply first.)

## Test Commands

```bash
# Terminal 1 — start API
go run ./cmd/api

# Terminal 2 — start worker
go run ./cmd/worker

# Watch checks appear
curl http://localhost:8080/api/endpoints/<id>/checks

# Ctrl+C worker — should stop cleanly (no panic)
```

## Definition of Done

- [ ] Worker runs as separate process
- [ ] Worker fetches endpoints from DB on each cycle
- [ ] Worker calls `CheckEndpoint` for each endpoint
- [ ] Results saved to DB
- [ ] Worker stops cleanly on Ctrl+C
- [ ] API still works while worker is running
- [ ] Can explain difference between goroutine and OS thread

## Questions Before Moving On

1. What is a goroutine? How is it different from a thread?
2. What is a channel used for?
3. What happens if you don't handle SIGTERM?
4. Why can spawning too many goroutines be a problem?

## Next: Module 5 — Scheduling
