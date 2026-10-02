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

Use the `go.opentelemetry.io/otel` packages. Create one package-level tracer (`otel.Tracer("beacon")`).

**Exercise: instrument a handler.** In a handler, at the top:

1. Start a span from the request's context: `ctx, span := tracer.Start(r.Context(), "handler.create_endpoint")`.
2. `defer span.End()` so the span closes when the handler returns.
3. Pass that `ctx` (not `r.Context()`) into the service call, so the child span links to this one.
4. On error, record it on the span (`span.RecordError(err)`) and set the span status to error.

**Exercise: thread it through the service.** The service's `Create` should take `ctx context.Context` as its **first** argument, start its own child span from that ctx, defer its end, and pass ctx down to the repository. Each layer starts a child span from the ctx it received — that's how the trace tree gets built.

## Context Propagation Rule

Every function that creates a span must accept `context.Context` as its **first** argument and pass it to child calls. A method like `Insert(ctx context.Context, e Endpoint)` keeps the trace connected; `Insert(e Endpoint)` (no ctx) breaks the chain — the child span becomes an orphan and the trace is lost. This is why idiomatic Go puts `ctx` first almost everywhere.

## Exporter Setup

Use Jaeger for local development. This is library wiring, not logic — read the OTel + Jaeger exporter docs and set up: a Jaeger exporter pointing at the local collector endpoint, a `TracerProvider` that batches to that exporter and tags itself with a service name, then register it globally with `otel.SetTracerProvider`. The skill here is *reading the library's setup docs and following them* — exactly what you'll do on the job.

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
