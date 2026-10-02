# Module 10 — Deploy Tracking

## Status: NOT STARTED

## What You're Building

Beacon becomes deploy-aware. When your app deploys, it tells Beacon. When an incident happens, Beacon can say: "This started 4 minutes after deploy #219."

## Why This Matters

Without deploy tracking:
> "API is DOWN."

With deploy tracking:
> "API went DOWN at 14:05. Deploy #219 happened at 14:00. Incident started 5 minutes after deploy."

This is the feature that makes Beacon more than a ping monitor.

## Concepts to Learn First

- **Event correlation** — linking two events by time proximity
- **Timestamps** — importance of accurate timezone-aware timestamps
- **Temporal queries** — SQL queries involving time ranges
- **Deployment metadata** — version, commit SHA, environment
- **GitHub Actions** — CI/CD pipeline basics, how to call a webhook from a pipeline

## Data Model

**Exercise: the `deployments` migration.** Write the `CREATE TABLE`. It needs:

- a UUID primary key
- a `project_id` UUID referencing `projects(id)`
- a `version` text column, not null (e.g. "219", "v1.4.2")
- a `commit` text column (git SHA)
- an `environment` text column defaulting to `'production'`
- a `deployed_at` timestamp defaulting to now
- a `deployed_by` text column

Question to answer as you write it: why `TIMESTAMPTZ` and not `TEXT` for `deployed_at`?

## New API Endpoint

```
POST /api/projects/:id/deployments   ← record a deployment
GET  /api/projects/:id/deployments   ← list deployments
```

### Request body
```json
{
  "version": "219",
  "commit": "abc123def456",
  "environment": "production",
  "deployed_by": "github-actions"
}
```

Protect this endpoint with an API key — only your CI pipeline should call it.

## GitHub Actions Integration

In your existing project's `.github/workflows/deploy.yml`:

```yaml
- name: Notify Beacon
  run: |
    curl -X POST https://your-beacon/api/projects/${{ env.PROJECT_ID }}/deployments \
      -H "Authorization: Bearer ${{ secrets.BEACON_API_KEY }}" \
      -H "Content-Type: application/json" \
      -d '{
        "version": "${{ github.run_number }}",
        "commit": "${{ github.sha }}",
        "environment": "production",
        "deployed_by": "${{ github.actor }}"
      }'
```

## Incident Correlation Query

When viewing an incident, find the most recent deploy before incident start.

**Exercise: the correlation query.** Write SQL that returns the single most recent deployment for a project (`$1`) whose `deployed_at` is *before* the incident's start time (`$2`). Hint: filter with `deployed_at < $2`, order newest-first, limit to one. (Also: name the columns explicitly instead of `SELECT *` — the habit you learned earlier.)

**Exercise: the time delta.** Given the incident and the deploy, compute how long after the deploy the incident started. Go's `time.Time` has a method that subtracts one time from another and returns a `time.Duration` — which method? The result prints as something like "4m32s".

## Dashboard Updates

### Endpoint detail page
```
Recent Incidents

DOWN at 14:05 — resolved at 14:23
  ↑ Deploy #219 was 5 minutes before this incident
```

### Deployment timeline
```
14:00  Deploy #219  (commit: abc123)
14:05  Incident OPEN
14:23  Incident RESOLVED
```

## File Structure Target

```
backend/
  internal/
    deployment/
      model.go        ← Deployment struct
      repository.go   ← DB queries
      service.go      ← correlation logic
      handler.go      ← API handlers
```

## API Key Middleware

**Exercise: `requireAPIKey`.** Middleware is a function that wraps a handler and returns a new handler. Write one that takes an `http.Handler` (call it `next`) and returns an `http.Handler`. In English:

1. Return an `http.HandlerFunc` (which adapts a function into a Handler).
2. Inside it, read the `Authorization` header from the request.
3. If it does not equal `"Bearer "` plus your `BEACON_API_KEY` env var, write a 401 (`http.Error` with `http.StatusUnauthorized`) and `return` — don't fall through.
4. Otherwise call `next.ServeHTTP(w, r)` to pass control to the real handler.

The shape — "function that takes a handler and returns a handler" — is the middleware pattern. Understand it once and you can write logging, auth, and rate-limit middleware the same way.

## Definition of Done

- [ ] `POST /api/projects/:id/deployments` stores deployment
- [ ] Endpoint protected by API key
- [ ] GitHub Actions (or curl) can record a deployment
- [ ] Incident detail shows nearest recent deployment
- [ ] Time delta calculated correctly
- [ ] Deployment timeline visible in dashboard
- [ ] Can explain event correlation

## Questions Before Moving On

1. Why store `deployed_at` as `TIMESTAMPTZ` not `TEXT`?
2. What SQL clause finds "most recent deployment before timestamp X"?
3. Why is the deployment endpoint protected by API key?
4. What is a GitHub Actions secret and how is it used?

## Next: Module 11 — Public Status Page
