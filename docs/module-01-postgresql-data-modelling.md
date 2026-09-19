# Module 1 — PostgreSQL & Data Modelling

## Status: DONE

## What You're Building

Replace in-memory slice from Module 0 with real PostgreSQL database. Data persists across restarts.

## Concepts to Learn First

- **Relational databases** — tables, rows, columns
- **Primary keys** — unique identifier per row
- **Foreign keys** — how tables reference each other
- **SQL** — SELECT, INSERT, UPDATE, DELETE
- **Indexes** — speed up queries
- **Transactions** — group operations, all succeed or all fail
- **Connection pools** — reuse open DB connections, don't open new one per request

**Recommended:** PostgreSQL official tutorial + `database/sql` Go docs.

## Data Model

### projects table
```sql
id          UUID PRIMARY KEY
name        TEXT NOT NULL
created_at  TIMESTAMPTZ DEFAULT NOW()
```

### endpoints table
```sql
id          UUID PRIMARY KEY
project_id  UUID REFERENCES projects(id)
name        TEXT NOT NULL
url         TEXT NOT NULL
interval    INTEGER NOT NULL  -- seconds
created_at  TIMESTAMPTZ DEFAULT NOW()
```

### checks table
```sql
id            UUID PRIMARY KEY
endpoint_id   UUID REFERENCES endpoints(id)
status        TEXT NOT NULL  -- 'up' | 'down'
status_code   INTEGER
response_time INTEGER        -- milliseconds
error         TEXT
checked_at    TIMESTAMPTZ DEFAULT NOW()
```

## Relationships

```
Project
  └── Endpoint (many per project)
        └── Check (many per endpoint)
```

## What to Build

1. Write migration SQL files in `backend/migrations/`
2. Connect Go to Postgres using `database/sql` + `lib/pq` driver
3. Replace in-memory slice with real DB queries
4. Implement repository layer — functions that talk to DB

## File Structure Target

```
backend/
  migrations/
    001_create_projects.sql
    002_create_endpoints.sql
    003_create_checks.sql
  internal/
    database/
      database.go     ← open connection, connection pool config
    project/
      repository.go   ← SQL queries for projects
      model.go        ← updated struct with UUID
```

## Repository Pattern

Handler should not write SQL. Repository does that.

```
Handler → calls Service → calls Repository → runs SQL → returns struct
```

**repository.go** example:
```go
func (r *Repository) List() ([]Project, error) {
    rows, err := r.db.Query("SELECT id, name, created_at FROM projects")
    // scan rows into structs
}

func (r *Repository) Create(name string) (Project, error) {
    // INSERT and return created row
}
```

## Environment Variables

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=beacon
DB_PASSWORD=secret
DB_NAME=beacon
```

Read with `os.Getenv()`. Never hardcode credentials.

## Test Commands

```bash
# Start postgres (docker)
docker run -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=beacon -p 5432:5432 postgres

# Run migrations manually
psql -h localhost -U postgres -d beacon -f migrations/001_create_projects.sql

# Test API
curl -X POST http://localhost:8080/api/projects -d '{"name":"Campus Exchange"}'
# Restart server — project should still be there
curl http://localhost:8080/api/projects
```

## Definition of Done

- [ ] Postgres running locally
- [ ] Migration files create all three tables
- [ ] Go connects to Postgres on startup
- [ ] `POST /api/projects` writes to DB
- [ ] `GET /api/projects` reads from DB
- [ ] Restart server — data persists
- [ ] Can explain why these are three separate tables

## Questions Before Moving On

1. Why use a foreign key instead of storing project name in the endpoint row?
2. What is a connection pool and why does it matter?
3. What happens if DB is down when server starts?
4. What is a transaction and when would you use one?

## Next: Module 2 — Endpoint Management API
