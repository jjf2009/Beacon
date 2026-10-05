package scheduler

import (
	"github.com/jjf2009/beacon/backend/internal/check"
)

type Scheduler struct {
    jobs map[string]*Job  // endpoint ID → job
    checkRepo   *check.Repository
}
