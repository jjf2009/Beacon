package stats

type Stats struct {
	UptimePercent   float64 `json:"uptime_percent"`
	AvgResponseTime int     `json:"avg_response_time"`
	TotalChecks     int     `json:"total_checks"`
	FailedChecks    int     `json:"failed_checks"`
}

type RawStats struct {
	TotalChecks     int
	UpChecks        int
	FailedChecks    int
	AvgResponseTime float64
}
