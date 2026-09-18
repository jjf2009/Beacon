package project

import (
	"database/sql"
	"fmt"
	"time"
)

type Repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) List() ([]Project, error){
	query := `SELECT id, name, created_at FROM projects`
	rows , err :=  r.db.Query(query)

	if err!= nil {
		return nil,err
	}
	defer rows.Close()
	projects := []Project{}
  for rows.Next() {
    var project Project

    if err := rows.Scan(&project.ID, &project.Name, &project.CreatedAt); err != nil {
      return nil, err
    }

   projects = append(projects, project)
  }

  if err := rows.Err(); err != nil {
    return nil, err
  }

  return projects, nil
}

func (r *Repository) Create(name string) (Project, error){
    query := `INSERT INTO projects (id, name, created_at) VALUES ($1, $2, $3) RETURNING id, name, created_at`

id := fmt.Sprintf("%d", time.Now().UnixNano())

var p Project
err := r.db.QueryRow(query, id, name, time.Now()).Scan(&p.ID, &p.Name, &p.CreatedAt)
if err != nil {
    return Project{}, err
}
return p, nil   // ← return the actual project, not empty Project{}
}


