package storage

import (
	"github.com/nabam/koel-server/models"
)

type Storage interface {
	CreateService(string) error
	DeleteService(string) error
	GetServices() ([]models.Service, error)
	Close() error
}
