package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/jjf2009/beacon/backend/internal/check"
	"github.com/jjf2009/beacon/backend/internal/checker"
	"github.com/jjf2009/beacon/backend/internal/endpoint"
)

type Worker struct {
	endpointRepo *endpoint.Repository
	checkRepo    *check.Repository
}

func New(endpointRepo *endpoint.Repository, checkRepo *check.Repository) *Worker {
	return &Worker{endpointRepo: endpointRepo, checkRepo: checkRepo}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	slog.Info("worker started, checking endpoints every 60s")

	for {
		select {
		case <-ctx.Done():
			// context cancelled — server is shutting down, exit cleanly
			slog.Info("worker stopping")
			return
		case <-ticker.C:
			// ticker fired — time to check all endpoints
			w.runChecks()
		}
	}
}

func (w *Worker) runChecks() {
	// fetch all registered endpoints from DB
	endpoints, err := w.endpointRepo.List()
	if err != nil {
		slog.Error("worker: failed to list endpoints", "error", err)
		return
	}

	slog.Info("worker: running checks", "count", len(endpoints))

	// check each endpoint and save the result
	for _, ep := range endpoints {
		result := checker.Check(ep.URL)

		if _, err := w.checkRepo.Save(ep.ID, result); err != nil {
			slog.Error("worker: failed to save check", "endpoint", ep.Name, "error", err)
			continue // don't stop — keep checking remaining endpoints
		}

		slog.Info("worker: check done", "endpoint", ep.Name, "status", result.Status, "latency_ms", result.ResponseTime)
	}
}
