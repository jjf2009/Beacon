package response

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "ok"
	StatusError = "error"
)

func WriteJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var msgs []string
	for _, e := range errs {
		switch e.ActualTag() {
		case "required":
			msgs = append(msgs, e.Field()+" is required")
		default:
			msgs = append(msgs, e.Field()+" is invalid")
		}
	}
	return Response{
		Status: StatusError,
		Error:  strings.Join(msgs, ", "),
	}
}
