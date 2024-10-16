package store

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	sbauth "github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
	"paws/internal/types"
)

type AuthStore interface {
	LogIn(email, password string) (sbauth.Session, error)
	SignUp(email, password, name string) (uuid.UUID, error)
	RefreshToken(token string) (sbauth.Session, error)
}

type UserStore interface {
	User(id uuid.UUID) (types.User, error)
	UserByEmail(email string) (types.User, error)
}

type PetStore interface {
	Pet(id uuid.UUID) (types.Pet, error)
	Pets(userID uuid.UUID) ([]types.Pet, error)
	Create(p *types.Pet) error
	Update(p *types.Pet, userID uuid.UUID) error
	Delete(id, userID uuid.UUID) error
}

type PostgresStore struct {
	AuthStore AuthStore
	PetStore  PetStore
	UserStore UserStore
}

func NewPostgresStore(db *sqlx.DB, sb *supabase.Client, jwtSecret string) *PostgresStore {
	return &PostgresStore{
		AuthStore: NewPostgresAuthStore(db, sb, jwtSecret),
		UserStore: NewPostgresUserStore(db),
		PetStore:  NewPostgresPetStore(db),
	}
}
