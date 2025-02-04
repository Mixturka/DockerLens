package dtos

type PingDTO struct {
	IP        string `json:"ip" validate:"required,ip"`
	IsSuccess bool   `json:"is_success"`
	Time      int64  `json:"ping_time" validate:"required"`
}
