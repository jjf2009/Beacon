package stats

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// GetEndpointStats queries checks from the last 24 hours for a given endpoint.
func (r *Repository) GetEndpointStats(endpointID string) (RawStats, error) {
	query := `
		SELECT
			COUNT(*) AS total_checks,
			COUNT(*) FILTER (WHERE status = 'up') AS up_checks,
			COUNT(*) FILTER (WHERE status != 'up') AS failed_checks,
			COALESCE(AVG(response_time), 0) AS avg_response_time
		FROM checks
		WHERE endpoint_id = $1
		  AND checked_at >= NOW() - INTERVAL '24 hours'`

	var raw RawStats
	err := r.db.QueryRow(query, endpointID).Scan(
		&raw.TotalChecks,
		&raw.UpChecks,
		&raw.FailedChecks,
		&raw.AvgResponseTime,
	)
	if err != nil {
		return RawStats{}, err
	}

	return raw, nil
}
