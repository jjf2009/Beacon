package project


type Service struct {
    repo *Repository   // ← add this
}

func NewService(repo *Repository) *Service {   // ← accept repo
    return &Service{repo: repo}
}

func (s *Service) Create(name string) (Project, error) {
    return s.repo.Create(name)   // ← delegate to repo
}

func (s *Service) List() ([]Project, error) {
    return s.repo.List()   // ← add this
}