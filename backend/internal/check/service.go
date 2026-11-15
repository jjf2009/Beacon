package check

import (
    "github.com/jjf2009/beacon/backend/internal/endpoint"
)

type Service struct {
    repo         *Repository
    endpointRepo *endpoint.Repository
}

func NewService(
    repo *Repository,
    endpointRepo *endpoint.Repository,
) *Service {
    return &Service{
        repo:         repo,
        endpointRepo: endpointRepo,
    }
}

func (s *Service) GetRecentChecks(
    endpointID string,
    limit int,
) ([]Check, error) {
    _, err := s.endpointRepo.GetById(endpointID)
    if err != nil {
        return nil, err
    }

    return s.repo.GetRecentChecks(endpointID, limit)
}