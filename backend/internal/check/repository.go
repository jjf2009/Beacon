package check

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jjf2009/beacon/backend/internal/checker"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Save stores a CheckResult from the checker package into the checks table.
func (r *Repository) Save(endpointID string, result checker.CheckResult) (Check, error) {
	query := `
		INSERT INTO checks (id, endpoint_id, status, status_code, response_time, error, checked_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id::text, endpoint_id::text, status, status_code, response_time, error, checked_at`

	id := uuid.New().String()
	var c Check
	err := r.db.QueryRow(query, id, endpointID, result.Status, result.StatusCode, result.ResponseTime, result.Error, time.Now()).
		Scan(&c.ID, &c.EndpointID, &c.Status, &c.StatusCode, &c.ResponseTime, &c.Error, &c.CheckedAt)
	if err != nil {
		return Check{}, err
	}
	return c, nil
}

// ListByEndpoint returns the most recent 20 checks for a given endpoint.
func (r *Repository) ListByEndpoint(endpointID string) ([]Check, error) {
	query := `
		SELECT id::text, endpoint_id::text, status, status_code, response_time, error, checked_at
		FROM checks
		WHERE endpoint_id = $1
		ORDER BY checked_at DESC
		LIMIT 20`

	rows, err := r.db.Query(query, endpointID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checks := []Check{}
	for rows.Next() {
		var c Check
		if err := rows.Scan(&c.ID, &c.EndpointID, &c.Status, &c.StatusCode, &c.ResponseTime, &c.Error, &c.CheckedAt); err != nil {
			return nil, err
		}
		checks = append(checks, c)
	}
	return checks, rows.Err()
}
