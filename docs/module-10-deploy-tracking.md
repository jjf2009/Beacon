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

### deployments table
```sql
id          UUID PRIMARY KEY
project_id  UUID REFERENCES projects(id)
version     TEXT NOT NULL      -- "219", "v1.4.2", "abc123"
commit      TEXT               -- git commit SHA
environment TEXT DEFAULT 'production'
deployed_at TIMESTAMPTZ DEFAULT NOW()
deployed_by TEXT               -- who triggered it
```

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

When viewing an incident, find the most recent deploy before incident start:

```sql
SELECT * FROM deployments
WHERE project_id = $1
  AND deployed_at < $2  -- before incident started_at
ORDER BY deployed_at DESC
LIMIT 1
```

Compute time delta:
```go
delta := incident.StartedAt.Sub(deploy.DeployedAt)
// "Incident started 4m32s after deploy #219"
```

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

Simple token check for deployment endpoint:

```go
func requireAPIKey(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        key := r.Header.Get("Authorization")
        if key != "Bearer "+os.Getenv("BEACON_API_KEY") {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

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
