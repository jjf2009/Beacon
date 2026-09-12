# Module 6 — Monitoring Dashboard

## Status: NOT STARTED

## What You're Building

Next.js dashboard that shows live monitoring data. First major milestone: Beacon can automatically monitor your existing applications.

## What the Dashboard Shows

### Endpoint status overview
```
Campus Exchange      ● UP     142ms    2 min ago
Resume Analyzer      ● UP     89ms     1 min ago
API Server           ● DOWN   —        30 sec ago
```

### Endpoint detail page
- Current status (UP/DOWN)
- Uptime percentage (last 24h)
- Average response time
- Last checked timestamp
- Recent checks table (last 20)
- Latency graph (last hour)

## New API Endpoints Needed

```
GET /api/endpoints/:id/checks          ← check history
GET /api/endpoints/:id/stats           ← uptime %, avg latency
GET /api/endpoints/:id/checks?limit=20 ← paginated checks
```

### Stats response shape
```json
{
  "uptime_percent": 99.2,
  "avg_response_time": 134,
  "total_checks": 1440,
  "failed_checks": 11
}
```

## Frontend Pages

```
/                        ← all endpoints, status overview
/endpoints/:id           ← detail page for one endpoint
/endpoints/new           ← add endpoint form (from Module 2)
```

## Key UI Concepts to Learn

- **Polling** — refresh data every 30s with `setInterval`
- **Loading states** — show spinner while fetching
- **Error states** — show message if fetch fails
- **Relative timestamps** — "2 min ago" not "2026-01-01T14:02:00Z"
- **Status colors** — green=UP, red=DOWN, grey=unknown

## Latency Graph

Use a simple charting library (`recharts` or `chart.js`). Plot last 60 check timestamps vs response times.

```
ms
200 │    ●
150 │  ●   ●   ●
100 │●           ●  ●
 50 │
    └──────────────────── time
```

## Polling vs Real-time

For now: poll every 30s with `setInterval`. Good enough.
Real-time (WebSockets) comes much later — don't add it now.

```typescript
useEffect(() => {
    fetchEndpoints()
    const interval = setInterval(fetchEndpoints, 30_000)
    return () => clearInterval(interval)  // cleanup on unmount
}, [])
```

## Uptime Calculation (backend)

```sql
SELECT
    COUNT(*) as total,
    COUNT(*) FILTER (WHERE status = 'up') as up_count
FROM checks
WHERE endpoint_id = $1
  AND checked_at > NOW() - INTERVAL '24 hours'
```

```go
uptime := float64(upCount) / float64(total) * 100
```

## Definition of Done

- [ ] Dashboard shows all endpoints with current status
- [ ] Status updates every 30s without page refresh
- [ ] Endpoint detail page shows check history
- [ ] Uptime % calculated correctly
- [ ] Latency graph renders
- [ ] DOWN endpoint shown in red
- [ ] Can navigate to detail page and back

## ── FIRST MAJOR MILESTONE ──

When this module is done, Beacon:
1. Accepts registered endpoints
2. Automatically checks them on schedule
3. Stores all results in Postgres
4. Displays live status in dashboard

You should be able to point Beacon at any of your existing projects and watch it monitor them.

## Questions Before Moving On

1. Why poll instead of WebSocket for now?
2. How do you calculate uptime percentage from raw check data?
3. What does `clearInterval` do and why is it important?
4. What SQL aggregation functions did you use?

## Next: Module 7 — Incident Detection
