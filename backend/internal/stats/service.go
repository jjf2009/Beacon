package stats

import (
	"math"

	"github.com/jjf2009/beacon/backend/internal/endpoint"
)

type Service struct {
	repo         *Repository
	endpointRepo *endpoint.Repository
}

func NewService(repo *Repository, endpointRepo *endpoint.Repository) *Service {
	return &Service{
		repo:         repo,
		endpointRepo: endpointRepo,
	}
}

func (s *Service) GetStats(endpointID string) (Stats, error) {
	// Verify that the endpoint exists
	_, err := s.endpointRepo.GetById(endpointID)
	if err != nil {
		return Stats{}, err
	}

	raw, err := s.repo.GetEndpointStats(endpointID)
	if err != nil {
		return Stats{}, err
	}

	if raw.TotalChecks == 0 {
		return Stats{
			UptimePercent:   0,
			AvgResponseTime: 0,
			TotalChecks:     0,
			FailedChecks:    0,
		}, nil
	}

	// Calculate uptime percentage: (up_checks / total_checks) * 100
	uptime := (float64(raw.UpChecks) / float64(raw.TotalChecks)) * 100
	// Round to 1 decimal place (e.g., 99.2)
	uptime = math.Round(uptime*10) / 10

	return Stats{
		UptimePercent:   uptime,
		AvgResponseTime: int(math.Round(raw.AvgResponseTime)),
		TotalChecks:     raw.TotalChecks,
		FailedChecks:    raw.FailedChecks,
	}, nil
}
