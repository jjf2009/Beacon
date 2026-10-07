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
    s.jobs[ep.ID] = &Job{endpoint: ep, ticker: ticker, stop: stop}
      
}


func (s *Scheduler) RemoveJob(id string){
    value,ok := s.jobs[id]
     if ok {
          close(value.stop)
          delete(s.jobs, id)
     }else {
          slog.Error("Scheduler does not exist", "id", id)
     }

}