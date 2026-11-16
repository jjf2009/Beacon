package check

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

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

func ListChecks(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting recent check results")

		// Extract endpoint ID
		id := r.PathValue("id")
		if id == "" {
			response.WriteJson(
				w,
				http.StatusBadRequest,
				response.GeneralError(fmt.Errorf("missing endpoint ID")),
			)
			return
		}

		// Default limit
		limit := 20

		// Parse custom limit if provided
		limitStr := r.URL.Query().Get("limit")
		if limitStr != "" {
			parsedLimit, err := strconv.Atoi(limitStr)
			if err != nil || parsedLimit <= 0 {
				response.WriteJson(
					w,
					http.StatusBadRequest,
					response.GeneralError(fmt.Errorf("limit must be a positive integer")),
				)
				return
			}
			limit = parsedLimit
		}

		// Fetch recent checks
		checks, err := svc.GetRecentChecks(id, limit)
		if err != nil {
			slog.Error("Failed to list checks", "endpoint_id", id, "error", err)
			response.WriteJson(
				w,
				http.StatusInternalServerError,
				response.GeneralError(err),
			)
			return
		}

		response.WriteJson(w, http.StatusOK, checks)
	}
}

