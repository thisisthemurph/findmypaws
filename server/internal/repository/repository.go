package repository

import (
	"errors"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type Repositories struct {
	AlertRepository AlertRepository
	PetRepository   PetRepository
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		AlertRepository: NewPostgresAlertRepository(db),
		PetRepository:   NewPostgresPetRepository(db),
	}
}
