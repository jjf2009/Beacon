package server
package server

import (
	"net/http"
)

func New() *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Beacon API is running"))
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}