# Module 14 — Grafana

## Status: NOT STARTED

## What You're Building

Connect Prometheus metrics to Grafana dashboards. Visualize Beacon's operational data with graphs real engineering teams use.

## What Grafana Is

Grafana = dashboard tool. Connects to data sources (Prometheus, Postgres, Loki, etc.) and renders graphs.

```
Prometheus (stores metrics)
      ↓
Grafana (visualizes metrics)
```

Grafana does not store data. It queries Prometheus on demand.

## Dashboards to Build

### 1. Beacon Operations Dashboard
- Checks per minute (rate)
- Check failure rate %
- P50 / P95 / P99 response time
- Active incidents (gauge)
- Endpoint count

### 2. Worker Dashboard
- Jobs processed per minute
- Jobs failed per minute
- Queue depth (if using Redis)
- Worker processing time

### 3. API Performance Dashboard
- HTTP requests per second
- Request latency histogram
- Error rate (4xx, 5xx)
- Endpoint-level breakdown

## Setup

### Docker Compose Addition

```yaml
grafana:
  image: grafana/grafana:latest
  ports:
    - "3001:3000"
  environment:
    - GF_SECURITY_ADMIN_PASSWORD=secret
  volumes:
    - grafana_data:/var/lib/grafana
```

Visit `http://localhost:3001`. Login: admin / secret.

### Add Prometheus Data Source

1. Grafana → Connections → Add data source
2. Select Prometheus
3. URL: `http://prometheus:9090`
4. Save & Test

## Example Panels

### Check Failure Rate
```
PromQL: rate(beacon_checks_total{status="down"}[5m]) / rate(beacon_checks_total[5m]) * 100
Panel: Stat or Time series
Unit: Percent (0-100)
Threshold: red at 5%
```

### P95 Response Time
```
PromQL: histogram_quantile(0.95, rate(beacon_check_duration_ms_bucket[5m]))
Panel: Time series
Unit: Milliseconds
```

### Active Incidents
```
PromQL: beacon_active_incidents
Panel: Stat
Thresholds: green=0, red>=1
```

## Dashboard as Code (JSON)

Grafana dashboards can be exported as JSON and version-controlled. Do this — don't recreate dashboards by hand after restart.

```
infra/
  grafana/
    dashboards/
      beacon-operations.json
      beacon-worker.json
```

## Alerting in Grafana

Grafana can send alerts when thresholds are crossed. For example:
- Check failure rate > 10% for 5 minutes → alert
- No checks recorded for 2 minutes → alert (worker down?)

This is separate from Beacon's own alerting (Module 8). Grafana alerts monitor Beacon itself. Beacon alerts monitor your applications.

## Docker Compose Full Stack

At this point `docker-compose.yml` should run everything:

```yaml
services:
  postgres:   ...
  redis:      ...
  api:        ...
  worker:     ...
  prometheus: ...
  grafana:    ...
  jaeger:     ...
```

Single command: `docker compose up` → full Beacon stack running.

## Definition of Done

- [ ] Grafana running on localhost:3001
- [ ] Prometheus connected as data source
- [ ] Operations dashboard with 4+ panels
- [ ] Check failure rate panel with threshold coloring
- [ ] P95 latency panel
- [ ] Active incidents stat panel
- [ ] Dashboard exported to JSON in repo
- [ ] Can explain: what Grafana does vs what Prometheus does

## Questions Before Moving On

1. What does Grafana actually store?
2. What is a PromQL `rate()` function doing?
3. What is a histogram quantile and what does P95 mean?
4. Why export dashboards as JSON?

## Next: Module 15 — Production Hardening
