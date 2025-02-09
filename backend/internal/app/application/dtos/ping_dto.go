package dtos

type PingDTO struct {
	IP        string `json:"ip" validate:"required,ip"`
	IsSuccess bool   `json:"is_success" validate:"required"`
	Duration  int64  `json:"ping_time" validate:"gte=0"`
}
