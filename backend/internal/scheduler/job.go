package scheduler

import (
	"time"

	"github.com/jjf2009/beacon/backend/internal/endpoint"
)

type Job struct {
    endpoint endpoint.Endpoint   // ← use the package prefix
    ticker   *time.Ticker
    stop     chan struct{}
}
