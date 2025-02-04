package entities

import "time"

type Ping struct {
	ID        string    `json:"id"`
	IP        string    `json:"ip"`
	IsSuccess bool      `json:"is_success"`
	Time      int64     `json:"ping_time"`
	CreatedAt time.Time `json:"time_stamp"`
}
