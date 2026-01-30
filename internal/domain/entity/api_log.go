package entity

import "time"

type APILog struct {
	ID           int64
	RequestID    string
	Timestamp    time.Time
	HTTPMethod   string
	RequestPath  string
	HTTPStatus   int
	LatencyMs    int64
	IPAddress    string
	UserAgent    string
	UserID       int64
	RequestBody  string
	ResponseBody string
	CreatedAt    time.Time
}
