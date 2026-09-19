package server

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/jjf2009/beacon/backend/internal/endpoint"
	"github.com/jjf2009/beacon/backend/internal/project"
)

func New(addr string, db *sql.DB) *http.Server {
	router := http.NewServeMux()
	repo :=project.NewRepository(db);
	service :=project.NewService(repo)
	erepo :=endpoint.NewRepository(db);
	eservice := endpoint.NewService(erepo);

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	router.HandleFunc("GET /api/projects", project.List(service))
	router.HandleFunc("POST /api/projects", project.Create(service))
	router.HandleFunc("POST /api/endpoints",endpoint.Create(eservice))
	router.HandleFunc("GET /api/endpoints",endpoint.List(eservice))
    router.HandleFunc("GET /api/endpoints/{id}", endpoint.GetById(eservice))
    router.HandleFunc("DELETE /api/endpoints/{id}", endpoint.Delete(eservice))


	return &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}
