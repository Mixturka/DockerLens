package interfaces

import "github.com/Mixturka/DockerLens/pinger/internal/domain/entities"

type ApiClient interface {
	SavePingResults(rs []entities.PingResult) error
}
