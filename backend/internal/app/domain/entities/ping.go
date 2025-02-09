package entities

import "time"

type Ping struct {
	IP          string
	IsSuccess   bool
	Duration    int64
	LastSuccess time.Time
}
