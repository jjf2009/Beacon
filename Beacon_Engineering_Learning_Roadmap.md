# Beacon — Engineering Learning & Build Roadmap

## Project Overview

**Beacon** is a deploy-aware monitoring and incident platform.

It monitors applications through health endpoints, detects failures, records check history, identifies incidents, and eventually correlates incidents with recent deployments.

A public status page is one component of Beacon — the core product is the monitoring and incident system behind it.

### Example

An application deploys version `#219`.

Four minutes later, its `/health` endpoint starts failing.

Beacon should eventually be able to show:

> **Incident detected — 4 minutes after deploy #219**

This project is intentionally designed as a learning project. The goal is not just to ship a portfolio project, but to understand the engineering concepts behind each part.

---

# Architecture

## Initial Architecture

```text
                    ┌─────────────────────┐
                    │     Next.js Web     │
                    │   Dashboard / UI    │
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │       Go API        │
                    │ REST API / Business │
                    │      Logic          │
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │    PostgreSQL       │
                    │   Persistent Data   │
                    └─────────────────────┘

                    ┌─────────────────────┐
                    │   Beacon Worker     │
                    │ Background Checks   │
                    └──────────┬──────────┘
                               │
                               ▼
                         HTTP Health
                            Checks
                               │
                               ▼
                    ┌─────────────────────┐
                    │ Existing Projects   │
                    │      /health        │
                    └─────────────────────┘
```

### Later Architecture

```text
Scheduler
    │
    ▼
Redis Queue
    │
    ├── Worker 1
    ├── Worker 2
    └── Worker 3
          │
          ▼
    HTTP Health Checks
          │
          ▼
      PostgreSQL
```

Redis, Prometheus, OpenTelemetry and Grafana are intentionally introduced later. Do not add infrastructure just because it is considered "industry standard."

---

# Learning Philosophy

Beacon should be built using a **learn → implement → test → explain** cycle.

For every module:

1. Learn the underlying concept.
2. Watch a tutorial or read documentation.
3. Understand why the technology exists.
4. Design the small feature yourself.
5. Implement it in Beacon.
6. Test it.
7. Ask:
   - Why does this technology exist?
   - What problem does it solve?
   - What would break if we removed it?
   - Can I explain this architecture without looking at the code?

## AI Usage Rule

AI should be used as a **learning accelerator**, not as a replacement for engineering.

### Good uses of AI

- Explain concepts.
- Explain unfamiliar code.
- Generate small examples.
- Help debug errors.
- Review your implementation.
- Suggest edge cases.
- Generate tests.
- Compare architectural alternatives.
- Explain documentation.

### Avoid

- Asking AI to build an entire module.
- Copying large generated code without understanding it.
- Adding technologies because AI recommends them.
- Building the project entirely through prompts.

### Recommended workflow

```text
Learn concept
      ↓
Write your own mini-spec
      ↓
Make architecture decisions
      ↓
Implement small piece
      ↓
Use AI when stuck
      ↓
Read and understand the code
      ↓
Test manually
      ↓
Ask AI to review
      ↓
Rebuild important parts yourself
```

---

# Beacon Module Roadmap

## Module 0 — Go Backend Foundation

### Learn

- Go syntax
- Packages
- Structs
- Interfaces
- Error handling
- Go modules
- Environment variables
- HTTP servers
- Routing
- JSON
- Middleware
- Basic backend project structure

### Build in Beacon

Create the first Go API.

Initial endpoints:

```text
GET  /health
GET  /api/projects
POST /api/projects
```

### Goal

Understand how:

```text
Browser
   ↓
HTTP Request
   ↓
Go Router
   ↓
Handler
   ↓
Business Logic
   ↓
Response
```

works.

---

# Module 1 — PostgreSQL & Data Modelling

## Learn

- Relational databases
- Tables
- Rows and columns
- Primary keys
- Foreign keys
- Relationships
- SQL
- CRUD
- Indexes
- Transactions
- Connection pools

## Initial Data Model

### Project

```text
id
name
created_at
```

### Endpoint

```text
id
project_id
name
url
interval
created_at
```

### Check

```text
id
endpoint_id
status
status_code
response_time
error
checked_at
```

## Relationship

```text
Project
   │
   └── Endpoint
          │
          └── Check
```

### Goal

Be able to explain why these are separate tables and how the relationships work.

---

# Module 2 — Endpoint Management API

Beacon needs to know **what it should monitor**.

## Features

```text
POST   /api/endpoints
GET    /api/endpoints
GET    /api/endpoints/:id
DELETE /api/endpoints/:id
```

## Example

```json
{
  "name": "Campus Exchange",
  "url": "https://example.com/health",
  "interval": 60
}
```

## Frontend

Create a simple Next.js dashboard where the user can:

- Add an endpoint.
- View registered endpoints.
- Delete an endpoint.

## Data Flow

```text
Next.js
   ↓
Go API
   ↓
PostgreSQL
```

### Goal

Understand the complete frontend → API → database flow.

---

# Module 3 — HTTP Monitoring

Now Beacon actually starts monitoring applications.

## Learn

- Go HTTP client
- HTTP status codes
- Request timeouts
- Latency measurement
- DNS failures
- Connection failures
- HTTP errors
- Context cancellation

## Beacon Behaviour

Given:

```text
https://myapp.com/health
```

Beacon sends an HTTP request.

It records:

```text
UP / DOWN
HTTP status code
Response time
Error message
Timestamp
```

### Example

```text
Endpoint: Campus Exchange
Status:   UP
Code:     200
Latency:  142ms
Checked:  21:30:12
```

### Goal

Be able to explain exactly what happens when Beacon checks an endpoint.

---

# Module 4 — Background Worker

The API should not be responsible for continuously checking endpoints.

Introduce a separate worker.

## Learn

- Goroutines
- Channels
- Concurrency
- Background processes
- Worker patterns
- Graceful shutdown

## Architecture

```text
Go API
  │
  └── Handles user requests

Go Worker
  │
  └── Performs monitoring checks
```

The worker:

1. Gets endpoints.
2. Checks them.
3. Records the result.
4. Repeats.

### Goal

Understand why background work should be separated from request/response handling.

---

# Module 5 — Scheduling

Beacon must check endpoints automatically.

## Learn

- Scheduling
- Timers
- Tickers
- Recurring jobs
- Job lifecycle

## Initial Behaviour

```text
Every 60 seconds
       ↓
Find endpoints
       ↓
Check endpoints
       ↓
Store results
```

Later:

```text
30 seconds
1 minute
5 minutes
```

depending on configuration.

### Goal

Beacon continuously monitors applications without requiring manual requests.

---

# Module 6 — Monitoring Dashboard

Now visualize the collected data.

## Dashboard should show

### Endpoint status

```text
Campus Exchange     UP
Resume Analyzer     UP
API Server          DOWN
```

### Endpoint details

- Current status
- Uptime
- Response time
- Last check
- Recent checks
- Latency graph

### Goal

Reach the first major milestone:

> **Beacon can automatically monitor my existing application.**

## First Milestone Flow

```text
Existing Project
      │
      │ /health
      ▼
    Beacon
      │
      ▼
  Go Worker
      │
      ▼
 HTTP Check
      │
      ▼
 PostgreSQL
      │
      ▼
 Next.js Dashboard
```

---

# Module 7 — Incident Detection

A single failed request does not necessarily mean an incident.

Beacon needs an incident model.

## Learn

- State machines
- Failure thresholds
- Deduplication
- Incident lifecycle

## Example Rule

```text
3 consecutive failures
        ↓
Create incident
```

## Initial Incident States

```text
OPEN
  ↓
INVESTIGATING
  ↓
RESOLVED
```

### Example

```text
Check 1 → UP
Check 2 → UP
Check 3 → DOWN
Check 4 → DOWN
Check 5 → DOWN
             ↓
       Incident Created
```

### Goal

Understand the difference between:

> "A request failed"

and:

> "The system is experiencing an incident."

---

# Module 8 — Alerting

Beacon should notify someone when an incident occurs.

## Initial Options

Start with one:

- Email
- Discord

## Flow

```text
Incident
   ↓
Alert
   ↓
Notification
```

## Learn

- Asynchronous notifications
- Retries
- Failure handling
- Notification policies

### Goal

Understand how systems communicate important events asynchronously.

---

# Module 9 — Redis & Queues

Only introduce Redis after the simple worker architecture is understood.

## Problem

Imagine Beacon monitors:

```text
10 endpoints
```

The simple worker is fine.

Now imagine:

```text
10,000 endpoints
```

One worker checking everything sequentially becomes a bottleneck.

## New Architecture

```text
Scheduler
    │
    ▼
Redis Queue
    │
    ├── Worker 1
    ├── Worker 2
    ├── Worker 3
    └── Worker N
```

## Learn

- Queues
- Producers
- Consumers
- Job retries
- Backoff
- Concurrency
- Distributed workers

### Goal

Understand **why** a queue exists instead of simply learning Redis as a technology.

---

# Module 10 — Deploy Tracking

This is where Beacon becomes deploy-aware.

## Store deployment events

Example:

```json
{
  "version": "219",
  "commit": "abc123",
  "environment": "production"
}
```

## Timeline

```text
14:00
Deploy #219

14:04
Latency increases

14:05
Health checks fail

14:06
Incident created
```

Beacon can then report:

> Incident started 5 minutes after deployment #219.

## Learn

- Event correlation
- Timestamps
- Temporal queries
- Deployment metadata
- GitHub Actions integration

### Goal

Move beyond basic monitoring into **deployment-aware incident analysis**.

---

# Module 11 — Public Status Page

Beacon should provide a public-facing status page.

## Example

```text
BEACON STATUS

API             ● Operational
Web Application  ● Operational
Database         ● Operational

Last updated: 30 seconds ago
```

## Learn

- Public/private data separation
- Public APIs
- Status aggregation
- Caching

### Goal

Allow users to communicate service health without exposing the private engineering dashboard.

---

# Module 12 — Prometheus

Now introduce metrics.

## Beacon metrics

Examples:

```text
beacon_checks_total
beacon_check_failures_total
beacon_check_duration
beacon_active_incidents
```

Prometheus periodically scrapes Beacon's `/metrics` endpoint.

## Learn

- Counters
- Gauges
- Histograms
- PromQL
- Metrics collection

### Goal

Understand how production systems expose measurable operational data.

---

# Module 13 — OpenTelemetry

Instrument Beacon itself.

## Trace example

```text
HTTP Request
     ↓
Go API
     ↓
Database Query
     ↓
Worker
     ↓
HTTP Health Check
```

Each part can become a trace/span.

## Learn

- Traces
- Spans
- Context propagation
- Instrumentation
- Exporters

### Goal

Understand what actually happens inside a distributed application request.

---

# Module 14 — Grafana

Connect:

```text
Prometheus
    ↓
Grafana
```

Create dashboards for:

- Endpoint latency
- Failure rate
- Check duration
- Worker jobs
- Active incidents
- API performance

### Goal

Build operational dashboards similar to what real engineering teams use.

---

# Module 15 — Production Hardening

Finally, make Beacon production-ready.

## Learn / Implement

- Docker
- Health checks
- Graceful shutdown
- Structured logging
- Configuration management
- Secrets
- CI/CD
- Integration tests
- Rate limiting
- Retries
- Monitoring Beacon itself
- Production deployment

### Goal

Deploy Beacon as a reliable service rather than just a local project.

---

# Recommended Build Order

```text
MODULE 0
Go Backend
    ↓
MODULE 1
PostgreSQL
    ↓
MODULE 2
Endpoint API
    ↓
MODULE 3
HTTP Monitoring
    ↓
MODULE 4
Background Worker
    ↓
MODULE 5
Scheduling
    ↓
MODULE 6
Monitoring Dashboard
    ↓
──── FIRST MAJOR MILESTONE ────
Beacon automatically monitors an application
    ↓
MODULE 7
Incident Detection
    ↓
MODULE 8
Alerting
    ↓
MODULE 9
Redis + Queues
    ↓
MODULE 10
Deploy Tracking
    ↓
MODULE 11
Public Status Page
    ↓
MODULE 12
Prometheus
    ↓
MODULE 13
OpenTelemetry
    ↓
MODULE 14
Grafana
    ↓
MODULE 15
Production Hardening
```

---

# Suggested Daily Learning Cycle

For each coding session:

## 1. Learn — 30–60 minutes

Watch a tutorial or read documentation about the specific concept you are about to use.

Do not try to learn the entire technology.

Learn only enough to understand the current module.

## 2. Design — 10–20 minutes

Write down:

- What am I building?
- What data do I need?
- What request/response should exist?
- Why am I using this technology?
- What alternatives exist?

## 3. Build — 2–3 hours

Implement the feature yourself.

Use AI only when necessary.

## 4. Test

Try:

- Normal input
- Invalid input
- Missing data
- Network failure
- Duplicate requests
- Unexpected responses

depending on the feature.

## 5. Explain

At the end of the session, answer:

```text
What did I learn today?

Why does this technology exist?

What part of Beacon uses it?

What would break if I removed it?

Can I explain the code without AI?
```

---

# Feature 1 — Start Here

The first practical feature should be:

## Register an Endpoint

Beacon should allow you to register an existing application's health endpoint.

### Example

```text
Name:
Campus Exchange

URL:
https://example.com/health

Interval:
60 seconds
```

### Backend

Create:

```text
POST /api/endpoints
GET  /api/endpoints
```

### Database

Store:

```text
id
name
url
interval
created_at
updated_at
```

### Frontend

Build a simple page:

```text
+-----------------------------+
| Add Endpoint                |
|                             |
| Name:     [_____________]   |
| URL:      [_____________]   |
| Interval: [60___________]   |
|                             |
|          [Add Endpoint]     |
+-----------------------------+

Registered Endpoints

Campus Exchange
https://example.com/health
```

### Definition of Done

You are finished when:

- [ ] You can add an endpoint.
- [ ] The endpoint is stored in PostgreSQL.
- [ ] Refreshing the page does not remove it.
- [ ] GET `/api/endpoints` returns it.
- [ ] You can see it in the dashboard.
- [ ] You understand the full request flow.

### Important

**Do not build the worker yet.**

The next feature will be:

> Beacon automatically checks the registered endpoint and stores the result.

---

# Engineering Principles for Beacon

## 1. Understand before abstracting

Do not create abstractions before you understand the problem.

## 2. Start simple

Use:

```text
Go
PostgreSQL
Next.js
```

before introducing:

```text
Redis
Prometheus
OpenTelemetry
Grafana
```

## 3. Add technology because of a problem

Not:

> "Real systems use Redis."

Instead:

> "Our workload now requires distributed job processing, so a queue solves this problem."

## 4. Build in layers

Every module should build on the previous one.

## 5. Prefer understanding over speed

A slower implementation that you understand is more valuable than a fast implementation generated by AI.

## 6. Keep an engineering journal

For every major feature, record:

```text
Feature:
What I built:

Concepts learned:

Why I chose this architecture:

What confused me:

What failed:

How I fixed it:

What happens at scale:

What I would change in production:
```

---

# Final Goal

By the end of Beacon, the goal is not simply to say:

> "I built a monitoring platform."

The goal is to be able to explain:

- How an HTTP server works.
- How a REST API is structured.
- How PostgreSQL stores related data.
- How background workers operate.
- How scheduled jobs work.
- How concurrent health checks work.
- How incidents are detected.
- How notifications are delivered.
- Why queues are useful.
- How deployment events can be correlated with failures.
- How metrics differ from logs and traces.
- How OpenTelemetry tracing works.
- How Prometheus collects metrics.
- How Grafana visualizes operational data.
- How to harden and deploy the system.

That is what turns Beacon from a **portfolio project** into an **engineering learning project**.
