package project

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jjf2009/beacon/backend/internal/utils/response"
)

type CreateRequest struct {
	Name string `json:"name" validate:"required"`
}

func List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projects := []Project{} // empty slice — encodes to [] not null
		response.WriteJson(w, http.StatusOK, projects)
	}
}

func Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating new project")

		var req CreateRequest

		// Decode request body into req struct
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Validate — checks `validate:"required"` tags on the struct
		if err := validator.New().Struct(req); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		// Build the project (in-memory for now, no DB yet)
		project := Project{
			ID:        "1",
			Name:      req.Name,
			CreatedAt: time.Now(),
		}

		slog.Info("project created", "name", project.Name)
		response.WriteJson(w, http.StatusCreated, project)
	}
}
