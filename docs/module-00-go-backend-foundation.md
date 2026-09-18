# Module 0 — Go Backend Foundation

## Status: DONE

## What You're Building

First Go HTTP server. No database yet. No business logic. Just: request comes in, Go handles it, response goes out.

## Concepts to Learn First

- **Go syntax** — variables, functions, structs, loops
- **Packages** — how Go organizes code (`package main`, `package server`)
- **Error handling** — Go returns `error` as a value, no exceptions
- **Go modules** — `go.mod`, `go get`, `import`
- **HTTP server** — how `net/http` works
- **Routing** — mapping URL paths to handler functions
- **JSON** — encoding/decoding with `encoding/json`
- **Middleware** — functions that wrap handlers

**Recommended:** tour.golang.org + read `net/http` package docs.

## What to Build

```
GET  /health           → { "status": "ok" }
GET  /api/projects     → [] (empty list, in-memory)
POST /api/projects     → create project (in-memory only)
```

No database. Store in Go slice. Data resets on restart — fine for now.

## File Structure Target

```
backend/
  cmd/api/main.go          ← starts server
  internal/
    server/server.go       ← route setup
    project/
      handler.go           ← HTTP handlers
      model.go             ← Project struct
```

## Request Flow

```
Browser → HTTP request
        → net/http receives
        → ServeMux matches route
        → Handler function runs
        → JSON response written
        → Browser receives
```

## Key Things to Implement

**model.go**
```go
type Project struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}
```

**server.go**
```go
mux := http.NewServeMux()
mux.HandleFunc("GET /health", healthHandler)
mux.HandleFunc("GET /api/projects", listProjects)
mux.HandleFunc("POST /api/projects", createProject)
```

## Test Commands

```bash
go run ./cmd/api

curl http://localhost:8080/health
curl http://localhost:8080/api/projects
curl -X POST http://localhost:8080/api/projects \
  -H "Content-Type: application/json" \
  -d '{"name": "Campus Exchange"}'
```

## Definition of Done

- [ ] Server starts on `:8080`
- [ ] `GET /health` returns `{"status":"ok"}`
- [ ] `POST /api/projects` stores project in memory
- [ ] `GET /api/projects` returns stored projects
- [ ] No crash on invalid JSON input
- [ ] Can explain `http.ServeMux` without looking it up

## Questions Before Moving On

1. What does `go run ./cmd/api` actually do?
2. Why does `ListenAndServe` block?
3. What is handler signature `(w http.ResponseWriter, r *http.Request)`?
4. Why is Go error handling a return value, not try/catch?

## Next: Module 1 — PostgreSQL
