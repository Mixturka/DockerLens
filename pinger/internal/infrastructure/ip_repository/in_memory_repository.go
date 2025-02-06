package iprepository

import (
	"errors"
	"net"
	"sync"
)

type InMemoryIpRepository struct {
	pingResults sync.Map
}

func NewInMemoryIpRepository() *InMemoryIpRepository {
	return &InMemoryIpRepository{
		pingResults: sync.Map{},
	}
}

func (r *InMemoryIpRepository) Add(ip string) error {
	parsedIp := net.ParseIP(ip)
	if parsedIp == nil {
		return errors.New("invalid IP")
	}
	r.pingResults.Store(parsedIp.String(), struct{}{})
	return nil
}

func (r *InMemoryIpRepository) Remove(ip string) error {
	r.pingResults.Delete(ip)
	return nil
}

func (r *InMemoryIpRepository) GetAll() ([]string, error) {
	var ips []string
	r.pingResults.Range(func(key, value interface{}) bool {
		ips = append(ips, key.(string))
		return true
	})
	return ips, nil
}
