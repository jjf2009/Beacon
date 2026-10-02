# Module 12 — Prometheus

## Status: NOT STARTED

## What You're Building

Expose operational metrics from Beacon itself. Prometheus scrapes a `/metrics` endpoint. You can then query how Beacon is performing.

## Why Metrics?

Logs tell you what happened. Metrics tell you how the system is behaving over time.

```
Log:    "check failed for endpoint abc123 at 14:05:32"
Metric: "beacon_check_failures_total = 142 (in last hour)"
```

Logs = events. Metrics = aggregates.

## Concepts to Learn First

- **Counters** — only go up, never reset (total checks run, total failures)
- **Gauges** — go up and down (active incidents, current endpoint count)
- **Histograms** — distribution of values (response time buckets)
- **PromQL** — Prometheus query language
- **Scraping** — Prometheus pulls metrics from `/metrics` endpoint on schedule

## Beacon Metrics to Expose

```
beacon_checks_total{status="up"}        counter  total UP checks
beacon_checks_total{status="down"}      counter  total DOWN checks
beacon_check_duration_ms               histogram  check response times
beacon_active_incidents                gauge      current open incidents
beacon_endpoints_total                 gauge      registered endpoint count
beacon_worker_jobs_processed_total     counter    queue jobs processed
beacon_worker_jobs_failed_total        counter    queue jobs failed
```

## Go Implementation

Use the official `prometheus/client_golang` library (import `prometheus` and `promhttp`). Read the library's docs — this is real DevOps skill: reading a library's API and wiring it in.

**Exercise: define the metrics.** As package-level variables:

1. A counter *vector* named `beacon_checks_total` with a label `"status"` (so you can count up vs down separately). Look up `prometheus.NewCounterVec` and `CounterOpts`.
2. A histogram named `beacon_check_duration_ms` with sensible buckets (e.g. 10, 50, 100, 200, 500, 1000, 5000 ms). Look up `prometheus.NewHistogram` and `HistogramOpts`.

**Exercise: register them.** Write an `init()` function that registers both metrics with Prometheus (`prometheus.MustRegister`). Why `init()`? What does Go do with it automatically?

**Exercise: expose `/metrics`.** Add one route to your mux that serves the Prometheus handler (`promhttp.Handler()`).

**Exercise: instrument the checker.** After a check runs, increment the counter for the result's status (`.WithLabelValues(status).Inc()`) and record the duration in the histogram (`.Observe(...)` — note it wants a `float64`, so convert). Where in your existing checker code does this go?

## Prometheus Config

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: beacon
    static_configs:
      - targets: ["localhost:8080"]
```

## Docker Compose Addition

```yaml
prometheus:
  image: prom/prometheus
  ports:
    - "9090:9090"
  volumes:
    - ./prometheus.yml:/etc/prometheus/prometheus.yml
```

## Useful PromQL Queries

```
# Check failure rate over 5 minutes
rate(beacon_checks_total{status="down"}[5m])

# 95th percentile response time
histogram_quantile(0.95, rate(beacon_check_duration_ms_bucket[5m]))

# Current open incidents
beacon_active_incidents
```

## File Structure Target

```
backend/
  internal/
    metrics/
      metrics.go    ← metric definitions + registration
```

## Definition of Done

- [ ] `/metrics` endpoint returns Prometheus format
- [ ] Prometheus scrapes Beacon every 15s
- [ ] Counter increments on every check
- [ ] Histogram records response times
- [ ] Gauge reflects current incident count
- [ ] Can write a PromQL query to find failure rate
- [ ] Can explain difference: counter vs gauge vs histogram

## Questions Before Moving On

1. Why is a counter different from a gauge?
2. What is a histogram bucket?
3. What does `rate()` do in PromQL?
4. Why does Prometheus pull metrics instead of having your app push them?

## Next: Module 13 — OpenTelemetry
