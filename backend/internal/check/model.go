package checker

import "time"

type CheckResult struct {
	Status string `json:"status"`
	StatusCode int`json:"statuscode"`
	ResponseTime int `json:"responsetime"`
	Error string `json:"error"`
    CheckedAt time.Time `json:"created_at"`
}