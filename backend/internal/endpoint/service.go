package endpoint

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(projectID string, name string, url string, interval int) (Endpoint, error) {
	return s.repo.Create(projectID, name, url, interval)
}

func (s *Service) List() ([]Endpoint, error) {
	return s.repo.List()
}

func (s *Service) GetById(id string) (Endpoint, error) {
	return s.repo.GetById(id)
}

func (s *Service) Delete(id string) (Endpoint, error) {
	return s.repo.Delete(id)
}
