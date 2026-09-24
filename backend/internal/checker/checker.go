package checker

import (
	"net/http"
	"time"
)

type CheckResult struct {
	Status       string    `json:"status"`
	StatusCode   int       `json:"status_code"`
	ResponseTime int64     `json:"response_time_ms"`
	Error        string    `json:"error,omitempty"`
	CheckedAt    time.Time `json:"checked_at"`
}

func Check(url string) CheckResult {
	// custom client with timeout — never use http.Get() directly,
	// it has no timeout and will hang forever if the server doesn't respond
	client := &http.Client{Timeout: 10 * time.Second}

	// record time before request so we can measure how long it takes
	start := time.Now()

	resp, err := client.Get(url)

	// calculate how many milliseconds the request took
	responseTime := time.Since(start).Milliseconds()

	// network-level failure: DNS error, connection refused, timeout, etc.
	// no HTTP response at all — status code is 0
	if err != nil {
		return CheckResult{
			Status:       "down",
			StatusCode:   0,
			ResponseTime: responseTime,
			Error:        err.Error(), // convert error interface to string
			CheckedAt:    time.Now(),
		}
	}
	defer resp.Body.Close()

	// HTTP error: server responded but with a failure status code (4xx, 5xx)
	if resp.StatusCode >= 400 {
		return CheckResult{
			Status:       "down",
			StatusCode:   resp.StatusCode,
			ResponseTime: responseTime,
			CheckedAt:    time.Now(),
		}
	}

	// success: server responded with 2xx or 3xx
	return CheckResult{
		Status:       "up",
		StatusCode:   resp.StatusCode,
		ResponseTime: responseTime,
		CheckedAt:    time.Now(),
	}
}
