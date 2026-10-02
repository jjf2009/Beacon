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

**Exercise: the `incidents` migration.** Write the `CREATE TABLE` yourself. The table needs:

- a UUID primary key
- an `endpoint_id` UUID that references `endpoints(id)`
- a `status` text column, not null (values: `'open'`, `'investigating'`, `'resolved'`)
- a `started_at` timestamp that defaults to now
- a `resolved_at` timestamp that is nullable (an open incident has no resolve time yet)

Model it on your earlier migrations. What column type did you use for the other timestamps?

## Detection Logic (runs after every check)

**Exercise: `EvaluateAfterCheck`.** Write a method on the incident service that takes an `endpointID string` and a check result, and returns an `error`. In plain English, the logic is:

- If the check status is `"down"`:
  - Count how many of the most recent checks failed in a row (a repo method).
  - If that count is at least 3, look up whether an open incident already exists for this endpoint.
  - Only create a new incident if there isn't one already (this is deduplication — don't open a second incident for the same outage).
- If the check status is `"up"`:
  - Look up any open incident for this endpoint, and if one exists, resolve it.

Think about: how do you represent "no open incident found"? What does a repo method return when nothing matches?

## Consecutive Failures Query

**Exercise: the query.** Write SQL that counts how many of the **3 most recent** checks for a given `endpoint_id` have status `'down'`. Hint: you need a subquery — first select the last 3 checks (`ORDER BY checked_at DESC LIMIT 3`), then count the down ones in the outer query. If the count is 3, all three most recent failed.

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
