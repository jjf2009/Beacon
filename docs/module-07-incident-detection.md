# Module 7 — Incident Detection

## Status: NOT STARTED

## What You're Building

Logic to decide: "is this endpoint experiencing an incident?" One failed check ≠ incident. Network blip happens. 3 consecutive failures = something is wrong.

## Core Concept

Difference between:
> "A request failed."

and:
> "The system is experiencing an incident."

## Concepts to Learn First

- **State machines** — system transitions between defined states
- **Failure thresholds** — N consecutive failures → trigger action
- **Deduplication** — don't create duplicate incidents for same outage
- **Incident lifecycle** — open → investigating → resolved

## Incident States

```
OPEN
  ↓ (engineer acknowledges)
INVESTIGATING
  ↓ (endpoint recovers)
RESOLVED
```

## Detection Rule

```
3 consecutive DOWN checks → create OPEN incident
1 UP check (while incident open) → auto-resolve to RESOLVED
```

Start simple. You can make rules configurable later.

## Data Model

### incidents table
```sql
id            UUID PRIMARY KEY
endpoint_id   UUID REFERENCES endpoints(id)
status        TEXT NOT NULL  -- 'open' | 'investigating' | 'resolved'
started_at    TIMESTAMPTZ DEFAULT NOW()
resolved_at   TIMESTAMPTZ
```

## Detection Logic (runs after every check)

```go
func (s *IncidentService) EvaluateAfterCheck(endpointID string, result CheckResult) error {
    if result.Status == "down" {
        consecutiveFails := repo.CountConsecutiveFails(endpointID)
        if consecutiveFails >= 3 {
            existing := repo.GetOpenIncident(endpointID)
            if existing == nil {
                repo.CreateIncident(endpointID)  // don't duplicate
            }
        }
    }

    if result.Status == "up" {
        existing := repo.GetOpenIncident(endpointID)
        if existing != nil {
            repo.ResolveIncident(existing.ID)
        }
    }
    return nil
}
```

## Consecutive Failures Query

```sql
SELECT COUNT(*) FROM (
    SELECT status FROM checks
    WHERE endpoint_id = $1
    ORDER BY checked_at DESC
    LIMIT 3
) recent
WHERE status = 'down'
```

If count = 3, all three most recent checks failed.

## New API Endpoints

```
GET  /api/incidents              ← all open incidents
GET  /api/incidents/:id          ← incident detail
PUT  /api/incidents/:id/status   ← update status (OPEN → INVESTIGATING)
GET  /api/endpoints/:id/incidents ← incidents for one endpoint
```

## Dashboard Updates

- Show open incident badge on DOWN endpoints
- Incident timeline on endpoint detail page
- Separate "Incidents" page listing all open incidents

## Example Timeline

```
Check 1 → UP
Check 2 → UP
Check 3 → DOWN    (1 fail)
Check 4 → DOWN    (2 fails)
Check 5 → DOWN    (3 fails) → Incident CREATED
Check 6 → DOWN    (incident already open, no duplicate)
Check 7 → UP      → Incident RESOLVED
```

## File Structure Target

```
backend/
  internal/
    incident/
      model.go        ← Incident struct
      repository.go   ← DB queries
      service.go      ← detection logic
      handler.go      ← API handlers
```

## Definition of Done

- [ ] 3 consecutive DOWN checks create one incident
- [ ] No duplicate incidents for same outage
- [ ] UP check resolves open incident
- [ ] `GET /api/incidents` returns open incidents
- [ ] Dashboard shows incident badge on affected endpoints
- [ ] Can manually move incident to INVESTIGATING
- [ ] Can explain the state machine without looking at code

## Questions Before Moving On

1. What is a state machine? Name the states and transitions.
2. Why 3 consecutive failures instead of 1?
3. What is deduplication and why does it matter here?
4. What SQL query finds the N most recent checks?

## Next: Module 8 — Alerting
