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

## Scheduler Design — write it yourself

### Exercise 1: the two structs

Define a `Scheduler` struct with:
- a field `jobs` that maps an endpoint ID (string) to a pointer-to-Job
- a field `checkRepo` that holds a pointer to the check package's Repository (the scheduler needs it to save results)

Define a `Job` struct with:
- a field holding the endpoint (type `endpoint.Endpoint`)
- a field holding a pointer to a `time.Ticker`
- a field `stop` that is a channel carrying an empty struct (a pure signal, no data)

Hints: map type syntax is `map[KeyType]ValueType`. An empty-struct channel is `chan struct{}`.

### Exercise 2: a constructor

Write a function `New` that takes a `*check.Repository` and returns a `*Scheduler`.
- Inside, build and return the address of a Scheduler whose `jobs` map is initialized (a nil map can't be written to — you must `make` it) and whose `checkRepo` is the one passed in.

### Exercise 3: AddJob

Write a method named `AddJob` on `*Scheduler` that takes one `endpoint.Endpoint` parameter (name it `ep`). Line by line, write Go that:

1. Creates a ticker that fires every `ep.Interval` seconds. (`time.NewTicker` wants a `time.Duration`, but `ep.Interval` is an `int` — you must convert it, then multiply by `time.Second`.)
2. Creates a `stop` channel that carries an empty struct.
3. Starts a **background goroutine** that loops forever. Inside the loop, wait on two channels at once (which keyword lets you wait on multiple channels?):
   - when the ticker's channel delivers a value: check `ep.URL`, then save the result using the scheduler's check repository, keyed by `ep.ID`
   - when the `stop` channel delivers: stop the ticker, then exit the goroutine
4. After the goroutine, store a pointer to a new `Job` (holding `ep`, the ticker, and the stop channel — use named fields) in the `jobs` map under key `ep.ID`.

### Exercise 4: RemoveJob

Write a method `RemoveJob` on `*Scheduler` that takes an `id string`. It should:
1. Look up the job in the map, using the comma-ok form to check it exists (`value, ok := m[key]`).
2. If it exists: close its `stop` channel (this wakes the goroutine's stop case), then delete the entry from the map with the builtin `delete`.

When you've written all four, ask me to check.

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
