# Module 2 — Endpoint Management API

## Status: DONE

## What You're Building

Beacon needs to know what to monitor. This module: full CRUD for endpoints + simple Next.js UI to manage them.

## What to Build

### Backend endpoints

```
POST   /api/endpoints        ← register new endpoint
GET    /api/endpoints        ← list all endpoints
GET    /api/endpoints/:id    ← get one endpoint
DELETE /api/endpoints/:id    ← remove endpoint
```

### Request body for POST

```json
{
  "name": "Campus Exchange",
  "url": "https://example.com/health",
  "interval": 60
}
```

### Frontend page

```
+-----------------------------+
| Add Endpoint                |
|                             |
| Name:     [_____________]   |
| URL:      [_____________]   |
| Interval: [60___________]   |
|                             |
|          [Add Endpoint]     |
+-----------------------------+

Registered Endpoints

Campus Exchange
https://example.com/health  — every 60s
[Delete]
```

## File Structure Target

```
backend/
  internal/endpoint/
    model.go        ← Endpoint struct
    repository.go   ← SQL queries
    service.go      ← business logic (validation)
    handler.go      ← HTTP handlers

frontend/
  app/
    page.tsx        ← endpoint list + add form
```

## Layer Responsibilities

```
Handler   → parse HTTP request, call service, write response
Service   → validate input, enforce rules, call repository
Repository → run SQL, return structs
```

## Validation Rules (in service layer)

- `name` must not be empty
- `url` must start with `http://` or `https://`
- `interval` must be >= 10 seconds

## Data Flow

```
Next.js form submit
      ↓
POST /api/endpoints
      ↓
Go handler parses body
      ↓
Service validates
      ↓
Repository inserts to Postgres
      ↓
Returns created endpoint
      ↓
Next.js shows in list
```

## Frontend Notes

- Use `fetch()` to call Go API
- `GET /api/endpoints` on page load
- `POST /api/endpoints` on form submit
- `DELETE /api/endpoints/:id` on delete button click
- No state management library needed — React `useState` + `useEffect` is enough

## CORS

Go API must allow requests from `http://localhost:3000`. Add middleware:

```go
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        // handle OPTIONS preflight
        next.ServeHTTP(w, r)
    })
}
```

## Test Commands

```bash
# Create endpoint
curl -X POST http://localhost:8080/api/endpoints \
  -H "Content-Type: application/json" \
  -d '{"name":"Campus Exchange","url":"https://example.com/health","interval":60}'

# List
curl http://localhost:8080/api/endpoints

# Get one
curl http://localhost:8080/api/endpoints/<id>

# Delete
curl -X DELETE http://localhost:8080/api/endpoints/<id>
```

## Definition of Done

- [ ] All four endpoints work
- [ ] Invalid input returns 400 with error message
- [ ] Frontend can add an endpoint
- [ ] Frontend lists all endpoints
- [ ] Frontend can delete an endpoint
- [ ] Page refresh does not lose data
- [ ] Can explain the full flow: form → API → DB → UI

## Questions Before Moving On

1. Why separate service layer from handler?
2. What HTTP status code for: not found / bad input / server error?
3. What is CORS and why does it exist?
4. What does `useEffect` do and when does it run?

## Next: Module 3 — HTTP Monitoring
