package server

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/jjf2009/beacon/backend/internal/check"
	"github.com/jjf2009/beacon/backend/internal/endpoint"
	"github.com/jjf2009/beacon/backend/internal/project"
)

func New(addr string, db *sql.DB) *http.Server {
	router := http.NewServeMux()

	// project routes
	prepo := project.NewRepository(db)
	psvc := project.NewService(prepo)
	router.HandleFunc("GET /api/projects", project.List(psvc))
	router.HandleFunc("POST /api/projects", project.Create(psvc))

	// endpoint routes
	erepo := endpoint.NewRepository(db)
	esvc := endpoint.NewService(erepo)
	router.HandleFunc("POST /api/endpoints", endpoint.Create(esvc))
	router.HandleFunc("GET /api/endpoints", endpoint.List(esvc))
	router.HandleFunc("GET /api/endpoints/{id}", endpoint.GetById(esvc))
	router.HandleFunc("DELETE /api/endpoints/{id}", endpoint.Delete(esvc))

	// check routes — run a manual check + view history
	crepo := check.NewRepository(db)
	router.HandleFunc("POST /api/endpoints/{id}/check", check.RunCheck(esvc, crepo))
	router.HandleFunc("GET /api/endpoints/{id}/checks", check.ListChecks(crepo))

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	return &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}
