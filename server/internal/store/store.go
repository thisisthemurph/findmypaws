package store

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/supabase-community/supabase-go"
	"paws/internal/types"
)

type UserStore interface {
	SignUp(email, password, name string) (types.User, error)
	User(id uuid.UUID) (types.User, error)
	UserByEmail(email string) (types.User, error)
}

type PetStore interface {
	Pet(id uuid.UUID) (types.Pet, error)
	Pets(userID uuid.UUID) ([]types.Pet, error)
	Create(p *types.Pet) error
	Update(p *types.Pet) error
}

type PostgresStore struct {
	PetStore  PetStore
	UserStore UserStore
}

func NewPostgresStore(db *sqlx.DB, sb *supabase.Client) *PostgresStore {
	return &PostgresStore{
		PetStore:  NewPostgresPetStore(db),
		UserStore: NewPostgresUserStore(db, sb),
	}
}
