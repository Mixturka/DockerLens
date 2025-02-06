package entities

import "time"

type Ping struct {
	ID        string
	IP        string
	IsSuccess bool
	Duration  int64
	CreatedAt time.Time
}
