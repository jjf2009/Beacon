package scheduler

import (
	"time"
	"github.com/jjf2009/beacon/backend/internal/endpoint"
)
type Job struct {
    endpoint endpoint.Endpoint   
    ticker   *time.Ticker
    stop     chan struct{}
}
