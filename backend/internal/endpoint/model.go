package endpoint

import "time"

type Endpoint struct {
    ID        string    `json:"id"`
    ProjectID string    `json:"project_id"`
    Name      string    `json:"name"`
    URL       string    `json:"url"`
    Interval  int       `json:"interval"`
    CreatedAt time.Time `json:"created_at"`
}
