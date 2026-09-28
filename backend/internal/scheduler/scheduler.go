package scheduler

import (
	"database/sql"
	"time"

	"github.com/jjf2009/beacon/backend/internal/checker"
	"github.com/jjf2009/beacon/backend/internal/endpoint"
)

type Scheduler struct {
    jobs map[string]*Job  // endpoint ID → job
    db   *sql.DB
}


func (s *Scheduler) AddJob(endpoint Endpoint) {
    ticker := time.NewTicker(time.Duration(endpoint.Interval) * time.Second)
    stop := make(chan struct{})
    
    go func() {
        for {
            select {
            case <-ticker.C:
                result := checker.Check(endpoint.URL)
                repo.SaveCheck(endpoint.ID, result)
            case <-stop:
                ticker.Stop()
                return
            }
        }
    }()
    
    s.jobs[endpoint.ID] = &Job{endpoint, ticker, stop}
}


func (s *Scheduler) RemoveJob(id string) {
    if job, ok := s.jobs[id]; ok {
        close(job.stop)
        delete(s.jobs, id)
    }
}