# Module 5 — Scheduling

## Status: NOT STARTED

## What You're Building

Worker from Module 4 checks endpoints on a fixed loop. This module: respect per-endpoint `interval` field. Campus Exchange checks every 60s. Critical API checks every 10s. Each on its own schedule.

## Concepts to Learn First

- **Timers** — fire once after duration
- **Tickers** — fire repeatedly at interval (`time.Ticker`)
- **Recurring jobs** — keep running on schedule
- **Job lifecycle** — start, pause, stop a scheduled job
- **Scheduler pattern** — manage many jobs with different intervals

## Problem with Simple Sleep Loop

Module 4 worker:
```
check all endpoints → sleep 60s → repeat
```

Problem: all endpoints checked at same interval. Endpoint with `interval=10` still waits 60s.

## What to Build

Scheduler that gives each endpoint its own ticker.

```
Endpoint A (interval=60)  → ticker fires every 60s → check A
Endpoint B (interval=30)  → ticker fires every 30s → check B
Endpoint C (interval=10)  → ticker fires every 10s → check C
```

## Scheduler Design

```go
type Scheduler struct {
    jobs map[string]*Job  // endpoint ID → job
    db   *sql.DB
}

type Job struct {
    endpoint Endpoint
    ticker   *time.Ticker
    stop     chan struct{}
}
```

### Start a job
```go
func (s *Scheduler) AddJob(endpoint Endpoint) {
    ticker := time.NewTicker(time.Duration(endpoint.Interval) * time.Second)
    stop := make(chan struct{})
    
    go func() {
        for {
            select {
            case <-ticker.C:
                result := checker.Check(endpoint.URL)
                repo.SaveCheck(endpoint.ID, result)
            case <-stop:
                ticker.Stop()
                return
            }
        }
    }()
    
    s.jobs[endpoint.ID] = &Job{endpoint, ticker, stop}
}
```

### Stop a job
```go
func (s *Scheduler) RemoveJob(id string) {
    if job, ok := s.jobs[id]; ok {
        close(job.stop)
        delete(s.jobs, id)
    }
}
```

## Dynamic Job Management

When endpoint is added via API → scheduler must start new job immediately.
When endpoint is deleted via API → scheduler must stop its job.

This requires communication between API process and Worker process. Simple approach: worker polls DB for changes every N seconds and reconciles jobs.

```
Worker loop:
1. Fetch all active endpoints from DB
2. Start jobs for endpoints not yet scheduled
3. Stop jobs for endpoints no longer in DB
4. Repeat every 30s (reconcile loop)
```

## File Structure Target

```
backend/
  internal/
    scheduler/
      scheduler.go   ← job management
      job.go         ← individual job struct
```

## Test Scenarios

- Add endpoint with `interval=10` → confirm checks appear every ~10s
- Add endpoint with `interval=60` → confirm slower check frequency
- Delete endpoint → confirm checks stop
- Restart worker → all jobs restart from DB

## Definition of Done

- [ ] Each endpoint checked at its own interval
- [ ] New endpoints picked up without restart
- [ ] Deleted endpoints stop being checked
- [ ] Worker handles 10+ endpoints concurrently without issue
- [ ] Graceful shutdown stops all tickers cleanly
- [ ] Can explain `time.Ticker` vs `time.Timer`

## Questions Before Moving On

1. What is the difference between `time.Ticker` and `time.Timer`?
2. Why use `select` with channels instead of `if` statements?
3. What happens to a goroutine if you forget to stop its ticker?
4. How would you handle an endpoint that takes 30s to respond when interval is 10s?

## Next: Module 6 — Monitoring Dashboard
