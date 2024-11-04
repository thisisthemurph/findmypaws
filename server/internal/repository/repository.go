package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFound = errors.New("not found")
)

type Repositories struct {
	PetRepository          PetRepository
	NotificationRepository NotificationRepository
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		PetRepository:          NewPetRepository(db),
		NotificationRepository: NewNotificationRepository(db),
	}
}
