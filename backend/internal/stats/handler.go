package stats

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jjf2009/beacon/backend/internal/utils/response"
)

func GetStats(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			response.WriteJson(
				w,
				http.StatusBadRequest,
				response.GeneralError(fmt.Errorf("missing endpoint ID")),
			)
			return
		}

		s, err := svc.GetStats(id)
		if err != nil {
			slog.Error("failed to get endpoint stats", "endpoint_id", id, "error", err)
			if strings.Contains(err.Error(), "not found") {
				response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
				return
			}
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, s)
	}
}
