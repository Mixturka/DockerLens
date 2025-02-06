package interfaces

type IpRepository interface {
	Add(ip string) error
	GetAll() ([]string, error)
	Remove(ip string) error
}
