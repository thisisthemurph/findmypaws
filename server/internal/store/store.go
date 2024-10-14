package store

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"paws/internal/types"
)

type PetStore interface {
	Pet(id uuid.UUID) (types.Pet, error)
	Pets(userID uuid.UUID) ([]types.Pet, error)
	Create(p *types.Pet) error
	Update(p *types.Pet) error
}

type PostgresStore struct {
	PetStore PetStore
}

func NewPostgresStore(db *sqlx.DB) *PostgresStore {
	return &PostgresStore{
		PetStore: NewPostgresPetStore(db),
	}
}
