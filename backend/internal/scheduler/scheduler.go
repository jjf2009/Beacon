package scheduler

import (
	"log/slog"
	"time"

	"github.com/jjf2009/beacon/backend/internal/check"
	"github.com/jjf2009/beacon/backend/internal/checker"
	"github.com/jjf2009/beacon/backend/internal/endpoint"
)

type Scheduler struct {
    jobs map[string]*Job  // endpoint ID → job
    checkRepo   *check.Repository
}


func New(repo *check.Repository) *Scheduler{
    return &Scheduler{jobs: make(map[string]*Job),checkRepo: repo}
}

func (s *Scheduler) AddJob(ep endpoint.Endpoint) {
    ticker := time.NewTicker(time.Duration(ep.Interval) * time.Second)

    stop := make(chan struct{})         // make a signal channel (issue #3)

    go func() {
        for {
            select {
            case <-ticker.C:
                res := checker.Check(ep.URL)        // call the real check fn on ep.URL (issue #4)
                if _, err := s.checkRepo.Save(ep.ID,res); err != nil {   // save with ep.ID (issue #5)
                    slog.Error("scheduler: save failed", "error", err)
                }
            case <-stop:
                ticker.Stop()      // stop the ticker HERE, not with defer (issue #1)
                return
            }
        }
    }()

    s.jobs[ep.ID] = &Job{endpoint: ep, ticker: ticker, stop: stop}
      // store the job in the map so RemoveJob can find it (issue #6)
}
