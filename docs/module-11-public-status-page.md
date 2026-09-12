# Module 11 — Public Status Page

## Status: NOT STARTED

## What You're Building

Public-facing page anyone can view — no login required. Shows current system status. Separate from internal engineering dashboard.

## Why Separate?

Internal dashboard: shows everything — raw checks, incidents, deploys, DB queries.
Public status page: shows only what users need to see — is the service up or down.

Never expose internal data to public. Separate route, separate data, separate concerns.

## Example Output

```
BEACON STATUS
─────────────────────────────────
API             ● Operational
Web Application  ● Operational  
Database         ● Degraded
─────────────────────────────────
Last updated: 30 seconds ago

Current Incidents
─────────────────
Database latency elevated — investigating
Started: 14 minutes ago
```

## Concepts to Learn First

- **Public/private data separation** — what to expose vs hide
- **Status aggregation** — compute overall system status from component statuses
- **Caching** — don't hit DB on every public request
- **Cache invalidation** — when to expire the cache

## Status Levels

```
Operational    → all checks UP, no incidents
Degraded       → some checks failing or incident INVESTIGATING
Outage         → incident OPEN
Unknown        → no recent checks
```

## New API Routes

```
GET /status                 ← public status page (Next.js route)
GET /api/public/status      ← JSON for status page (no auth)
```

Public API returns only safe data:

```json
{
  "status": "degraded",
  "components": [
    { "name": "API", "status": "operational" },
    { "name": "Web App", "status": "operational" },
    { "name": "Database", "status": "degraded" }
  ],
  "active_incidents": [
    {
      "title": "Database latency elevated",
      "status": "investigating",
      "started_at": "2026-01-01T14:00:00Z"
    }
  ],
  "updated_at": "2026-01-01T14:14:30Z"
}
```

Note: no IDs, no URLs, no internal details in public response.

## Caching

Public status page can get many requests. Don't query DB every time.

Simple cache: store computed status in memory, expire every 30s.

```go
type StatusCache struct {
    mu        sync.Mutex
    data      *PublicStatus
    expiresAt time.Time
}

func (c *StatusCache) Get() (*PublicStatus, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    if time.Now().Before(c.expiresAt) {
        return c.data, true
    }
    return nil, false
}
```

Later you could use Redis for cache if you want it to survive restarts.

## Frontend Route

```
/status    ← public, no auth, anyone can access
```

Different layout from internal dashboard. Minimal. Clean.

## 90-Day Uptime History

Many status pages show uptime history per component:

```
API  ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░▓▓▓▓▓▓▓▓  99.1%
      ← 90 days →
```

Each bar = one day. Green = operational. Red = outage. Query checks table for daily uptime.

## File Structure Target

```
backend/
  internal/
    status/
      service.go      ← aggregate status from checks + incidents
      handler.go      ← public API handler

frontend/
  app/
    status/
      page.tsx        ← public status page
```

## Definition of Done

- [ ] `/status` page accessible without login
- [ ] Shows current status per component
- [ ] Shows active incidents (without internal details)
- [ ] Status cached — DB not hit on every request
- [ ] No internal data (IDs, URLs, raw errors) exposed publicly
- [ ] 90-day uptime bars rendered
- [ ] Can explain difference between public and private API

## Questions Before Moving On

1. What data should never appear on a public status page?
2. Why cache the status response?
3. What is cache invalidation and why is it hard?
4. How do you aggregate "overall status" from multiple component statuses?

## Next: Module 12 — Prometheus
