package endpoint

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jjf2009/beacon/backend/internal/utils/response"
)

type CreateRequest struct {
	ProjectID string `json:"project_id" validate:"required"`
	Name      string `json:"name"       validate:"required"`
	URL       string `json:"url"        validate:"required,url"`
	Interval  int    `json:"interval"   validate:"required,min=10"`
}

func List(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoints, err := svc.List()
		if err != nil {
			slog.Error("error listing endpoints", "error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, endpoints)
	}
}

func Create(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating new endpoint")

		var req CreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		ep, err := svc.Create(req.ProjectID, req.Name, req.URL, req.Interval)
		if err != nil {
			slog.Error("error creating endpoint", "error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("endpoint created", "name", ep.Name)
		response.WriteJson(w, http.StatusCreated, ep)
	}
}

func GetById(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("missing id")))
			return // ← was missing
		}

		ep, err := svc.GetById(id)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
			return // ← was missing
		}
		response.WriteJson(w, http.StatusOK, ep)
	}
}

func Delete(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("missing id")))
			return // ← was missing
		}

		ep, err := svc.Delete(id)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
			return // ← was missing
		}
		response.WriteJson(w, http.StatusOK, ep)
	}
}
