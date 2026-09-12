# Module 13 — OpenTelemetry

## Status: NOT STARTED

## What You're Building

Instrument Beacon with distributed tracing. See exactly what happens inside a request — which DB query was slow, where time was spent, what called what.

## Why Tracing?

Metrics tell you there's a problem. Logs tell you what happened. Traces tell you where time went.

```
Incoming request: POST /api/endpoints
  └─ handler.go                    2ms
      └─ service.Validate()        0.3ms
      └─ repository.Insert()       18ms
          └─ sql: INSERT ...       17ms  ← slow!
```

Without traces: "requests are slow." With traces: "the DB insert takes 17ms."

## Concepts to Learn First

- **Traces** — end-to-end record of one request's path through the system
- **Spans** — one unit of work within a trace (one function call, one DB query)
- **Context propagation** — pass trace context through function calls and HTTP requests
- **Instrumentation** — adding trace code to your app
- **Exporters** — where traces are sent (Jaeger, Tempo, Honeycomb)

## Trace Structure

```
Trace: "POST /api/endpoints"  (total: 21ms)
├── span: http.handler         (21ms)
│   ├── span: service.validate  (0.3ms)
│   └── span: repo.insert       (18ms)
│       └── span: db.exec       (17ms)
```

## OpenTelemetry vs Prometheus

| | Prometheus | OpenTelemetry |
|---|---|---|
| What | Metrics (numbers) | Traces (request paths) |
| When | Aggregate view | Per-request detail |
| Question answered | "How many failures/hour?" | "Why was this request slow?" |

## Go Implementation

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("beacon")

func (h *Handler) CreateEndpoint(w http.ResponseWriter, r *http.Request) {
    ctx, span := tracer.Start(r.Context(), "handler.create_endpoint")
    defer span.End()

    // pass ctx through all calls
    endpoint, err := h.service.Create(ctx, req)
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
    }
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (Endpoint, error) {
    ctx, span := tracer.Start(ctx, "service.create")
    defer span.End()

    return s.repo.Insert(ctx, req)
}
```

## Context Propagation Rule

Every function that creates a span must accept `context.Context` as first argument and pass it to child calls.

```go
// correct
func (r *Repo) Insert(ctx context.Context, e Endpoint) error

// wrong — trace context lost
func (r *Repo) Insert(e Endpoint) error
```

## Exporter Setup

Use Jaeger for local development:

```go
exporter, _ := jaeger.New(
    jaeger.WithCollectorEndpoint(
        jaeger.WithEndpoint("http://localhost:14268/api/traces"),
    ),
)

tp := sdktrace.NewTracerProvider(
    sdktrace.WithBatcher(exporter),
    sdktrace.WithResource(resource.NewWithAttributes(
        semconv.ServiceNameKey.String("beacon-api"),
    )),
)
otel.SetTracerProvider(tp)
```

## Docker Compose Addition

```yaml
jaeger:
  image: jaegertracing/all-in-one:latest
  ports:
    - "16686:16686"   # Jaeger UI
    - "14268:14268"   # collector
```

Visit `http://localhost:16686` to view traces.

## What to Instrument

Priority order:
1. HTTP handlers (every incoming request)
2. DB queries (usually the bottleneck)
3. HTTP health checks (outgoing requests)
4. Queue job processing (worker)

## File Structure Target

```
backend/
  internal/
    telemetry/
      tracer.go    ← OTel setup, tracer init
```

## Definition of Done

- [ ] Jaeger running locally
- [ ] API traces visible in Jaeger UI
- [ ] Each handler creates a span
- [ ] DB calls appear as child spans
- [ ] Errors recorded on spans
- [ ] Context passed through all instrumented functions
- [ ] Can explain: trace vs span vs context propagation

## Questions Before Moving On

1. What is a span? What is a trace?
2. Why must `context.Context` be the first argument?
3. What is the difference between a trace and a log?
4. What is "context propagation" across service boundaries?

## Next: Module 14 — Grafana
