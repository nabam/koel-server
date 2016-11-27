package models

import (
	"github.com/satori/go.uuid"
)

type Service struct {
	Id   uuid.UUID
	Name string
}

func NewService(name string) Service {
	return Service{Id: uuid.NewV4(), Name: name}
}
