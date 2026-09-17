package project

import "time"

type Service struct {
    // empty for now — Module 1 adds db *sql.DB here
}

func NewService() *Service {
    return &Service{}
}

func (s *Service) Create(name string) (Project, error) {
		project := Project{
			ID:        "1",
			Name:      name,
			CreatedAt: time.Now(),
		}

		return project,nil
}
	