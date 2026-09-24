package check

import (
	"log/slog"
	"net/http"

	"github.com/jjf2009/beacon/backend/internal/checker"
	"github.com/jjf2009/beacon/backend/internal/endpoint"
	"github.com/jjf2009/beacon/backend/internal/utils/response"
)

// RunCheck manually triggers a health check for one endpoint and saves the result.
func RunCheck(esvc *endpoint.Service, repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		// fetch endpoint so we have its URL
		ep, err := esvc.GetById(id)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
			return
		}

		// run the actual HTTP check
		result := checker.Check(ep.URL)

		// save result to DB
		saved, err := repo.Save(id, result)
		if err != nil {
			slog.Error("failed to save check result", "error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("check completed", "endpoint", ep.Name, "status", result.Status, "latency_ms", result.ResponseTime)
		response.WriteJson(w, http.StatusOK, saved)
	}
}

// ListChecks returns the recent check history for one endpoint.
func ListChecks(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		checks, err := repo.ListByEndpoint(id)
		if err != nil {
			slog.Error("failed to list checks", "error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, checks)
	}
}
