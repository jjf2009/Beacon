package scheduler

import (
	"log/slog"
	"sync"
	"time"

	"github.com/jjf2009/beacon/backend/internal/check"
	"github.com/jjf2009/beacon/backend/internal/checker"
	"github.com/jjf2009/beacon/backend/internal/endpoint"
)

type Scheduler struct {
    mu        sync.Mutex
    jobs      map[string]*Job
    checkRepo *check.Repository
}


func New(repo *check.Repository) *Scheduler{
    return &Scheduler{jobs: make(map[string]*Job),checkRepo: repo}
}

func (s *Scheduler) AddJob(ep endpoint.Endpoint) {
    ticker := time.NewTicker(time.Duration(ep.Interval) * time.Second)
    stop := make(chan struct{})         
    go func() {
        for {
            select {
            case <-ticker.C:
                res := checker.Check(ep.URL)        
                if _, err := s.checkRepo.Save(ep.ID,res); err != nil {   
                    slog.Error("scheduler: save failed", "error", err)
                }
            case <-stop:
                ticker.Stop()      
                return
            }
        }
    }()
    if _, exists := s.jobs[ep.ID]; exists {
    // Already scheduled
    return}
    s.mu.Lock()
    s.jobs[ep.ID] = &Job{endpoint: ep, ticker: ticker, stop: stop}
    s.mu.Unlock()
      
}


func (s *Scheduler) RemoveJob(id string){
    value,ok := s.jobs[id]
     if ok {
          close(value.stop)
          s.mu.Lock()
          delete(s.jobs, id)
          s.mu.Unlock()
     }else {
          slog.Error("Scheduler does not exist", "id", id)
     }

}