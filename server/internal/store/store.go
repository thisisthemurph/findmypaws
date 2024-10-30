package store

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"paws/internal/types"
)

type PetStore interface {
	Pet(id uuid.UUID) (types.Pet, error)
	Pets(userID string) ([]types.Pet, error)
	Create(p *types.Pet) error
	Update(p *types.Pet, userID string) error
	UpdateAvatar(avatarURI string, id uuid.UUID, userID string) error
	Delete(id uuid.UUID, userID string) error
}

type AlertStore interface {
	Create(alert types.Alert) error
}

type PostgresStore struct {
	AlertStore AlertStore
	PetStore   PetStore
}

func NewPostgresStore(db *sqlx.DB) *PostgresStore {
	return &PostgresStore{
		AlertStore: NewPostgresAlertStore(db),
		PetStore:   NewPostgresPetStore(db),
	}
}
