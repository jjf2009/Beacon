package check

import "time"

// Check represents a stored result of one health check in the database.
// Separate from checker.CheckResult — that's the live result,
// this is the persisted record with its own DB id and endpoint reference.
type Check struct {
	ID           string    `json:"id"`
	EndpointID   string    `json:"endpoint_id"`
	Status       string    `json:"status"`
	StatusCode   int       `json:"status_code"`
	ResponseTime int64     `json:"response_time_ms"`
	Error        string    `json:"error,omitempty"`
	CheckedAt    time.Time `json:"checked_at"`
}
