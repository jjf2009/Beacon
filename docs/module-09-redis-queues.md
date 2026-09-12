# Module 9 — Redis & Queues

## Status: NOT STARTED

## What You're Building

Replace the simple worker loop with a queue-based architecture. Introduce Redis. Understand why queues exist before adding them.

## Why a Queue?

### Current system (works fine for small scale)
```
Worker → fetch all endpoints → check one by one → sleep → repeat
```

Scales to maybe ~100 endpoints on one machine. Good enough to start.

### Problem at scale
```
10,000 endpoints
  × 60s interval
  = 167 checks/second
```

One sequential worker cannot keep up. Checks fall behind schedule.

### Queue solution
```
Scheduler (1) → enqueues check jobs → Redis Queue
                                          ↓
                              Worker 1 ─────┤
                              Worker 2 ─────┤ pull jobs, run checks
                              Worker N ─────┘
```

Add more workers → horizontal scaling.

## Concepts to Learn First

- **Queues** — FIFO data structure, producer adds, consumer removes
- **Producers** — put jobs into queue (scheduler)
- **Consumers** — take jobs from queue (workers)
- **Job retries** — if check fails, re-enqueue with attempt count
- **Backoff** — wait longer between retries (1s, 2s, 4s, 8s...)
- **Dead letter queue** — jobs that failed too many times go here
- **Distributed workers** — multiple processes consuming same queue

## Redis Basics

Redis = in-memory data store. Fast. Not a database replacement.

Key commands used for queues:
```
LPUSH queue job   ← add job to left of list
RPOP  queue       ← remove job from right (FIFO)
BRPOP queue 30    ← blocking pop, wait up to 30s for job
```

## Job Shape

```go
type CheckJob struct {
    EndpointID string    `json:"endpoint_id"`
    URL        string    `json:"url"`
    AttemptNum int       `json:"attempt_num"`
    EnqueuedAt time.Time `json:"enqueued_at"`
}
```

## Architecture

```
Scheduler
  - every N seconds, fetch due endpoints from DB
  - LPUSH check job to Redis queue

Worker(s)
  - BRPOP from Redis queue (blocking, efficient)
  - run CheckEndpoint(url)
  - save result to Postgres
  - if fail, re-enqueue with attempt_num++
  - if attempt_num > 3, discard (or dead-letter)
```

## File Structure Target

```
backend/
  internal/
    queue/
      queue.go      ← Redis queue interface
      redis.go      ← Redis implementation
    scheduler/
      scheduler.go  ← updated: enqueues jobs instead of direct check
  cmd/
    worker/main.go  ← updated: pulls from queue
```

## Docker Compose Addition

```yaml
redis:
  image: redis:7
  ports:
    - "5379:6379"
```

## Test Scenario

1. Start Redis
2. Start scheduler (enqueues jobs)
3. Start 3 worker instances
4. Watch all three consume jobs concurrently
5. Kill one worker → others keep processing
6. Restart killed worker → picks up new jobs

## Definition of Done

- [ ] Scheduler enqueues check jobs to Redis
- [ ] Workers consume from Redis queue
- [ ] Multiple workers can run simultaneously
- [ ] Failed jobs retried up to 3 times
- [ ] Jobs not lost if one worker crashes
- [ ] Can run with 2+ worker instances
- [ ] Can explain producer/consumer pattern

## Questions Before Moving On

1. What is FIFO? Why does order matter for check jobs?
2. What is `BRPOP` and why is it better than polling?
3. What happens to a job if the worker crashes mid-check?
4. Why not use Postgres as a queue instead of Redis?

## Next: Module 10 — Deploy Tracking
