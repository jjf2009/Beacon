package server

import (
	"net/http"
	"time"

	"github.com/jjf2009/beacon/backend/internal/project"
)

func New(addr string) *http.Server {
	router := http.NewServeMux()
	service :=project.NewService()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	router.HandleFunc("GET /api/projects", project.List(service))
	router.HandleFunc("POST /api/projects", project.Create(service))

	return &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}
