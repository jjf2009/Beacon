package endpoint

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List() ([]Endpoint, error) {
	query := `SELECT id::text, project_id::text, name, url, interval, created_at FROM endpoints`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	endpoints := []Endpoint{}
	for rows.Next() {
		var ep Endpoint
		if err := rows.Scan(&ep.ID, &ep.ProjectID, &ep.Name, &ep.URL, &ep.Interval, &ep.CreatedAt); err != nil {
			return nil, err
		}
		endpoints = append(endpoints, ep)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return endpoints, nil
}

func (r *Repository) Create(projectID string, name string, url string, interval int) (Endpoint, error) {
	query := `
		INSERT INTO endpoints (id, project_id, name, url, interval, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id::text, project_id::text, name, url, interval, created_at`

	id := uuid.New().String()
	var ep Endpoint
	err := r.db.QueryRow(query, id, projectID, name, url, interval, time.Now()).
		Scan(&ep.ID, &ep.ProjectID, &ep.Name, &ep.URL, &ep.Interval, &ep.CreatedAt)
	if err != nil {
		return Endpoint{}, err
	}
	return ep, nil
}

func (r *Repository) GetById(id string) (Endpoint, error) {
	query := `SELECT id::text, project_id::text, name, url, interval, created_at FROM endpoints WHERE id = $1`

	var ep Endpoint
	err := r.db.QueryRow(query, id).Scan(&ep.ID, &ep.ProjectID, &ep.Name, &ep.URL, &ep.Interval, &ep.CreatedAt)
	if err == sql.ErrNoRows {
		return Endpoint{}, fmt.Errorf("endpoint %s not found", id)
	}
	if err != nil {
		return Endpoint{}, err
	}
	return ep, nil
}

func (r *Repository) Delete(id string) (Endpoint, error) {
	query := `
		DELETE FROM endpoints WHERE id = $1
		RETURNING id::text, project_id::text, name, url, interval, created_at`

	var ep Endpoint
	err := r.db.QueryRow(query, id).Scan(&ep.ID, &ep.ProjectID, &ep.Name, &ep.URL, &ep.Interval, &ep.CreatedAt)
	if err == sql.ErrNoRows {
		return Endpoint{}, fmt.Errorf("endpoint %s not found", id)
	}
	if err != nil {
		return Endpoint{}, err
	}
	return ep, nil
}
