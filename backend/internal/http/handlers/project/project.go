package project

import (
	"net/http"

	"github.com/jjf2009/beacon/backend/internal/untils/response"
)



func Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		response.WriteJson(w, http.StatusOK)
	}
}
