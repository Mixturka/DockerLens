package entities

import "time"

type Ping struct {
	ID        string
	IP        string
	IsSuccess bool
	Time      int64
	CreatedAt time.Time
}
