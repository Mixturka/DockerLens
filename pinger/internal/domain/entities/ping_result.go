package entities

type PingResult struct {
	IP        string `json:"ip"`
	IsSuccess bool   `json:"is_success"`
	Duration  int64  `json:"ping_time"`
}
