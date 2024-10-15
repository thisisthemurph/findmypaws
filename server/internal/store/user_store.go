package store

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"paws/internal/types"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUserWithEmailExists = errors.New("email already exists")
)

type PostgresUserStore struct {
	*sqlx.DB
}

func NewPostgresUserStore(db *sqlx.DB) *PostgresUserStore {
	return &PostgresUserStore{db}
}

func (s *PostgresUserStore) User(id uuid.UUID) (types.User, error) {
	stmt := `
		select p.*, u.id as auth_id, u.email, u.created_at, u.last_sign_in_at
		from public.profiles p
		join auth.users u on u.id = p.id
		where p.id = $1;`

	var u types.User
	if err := s.Get(&u, stmt, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.User{}, ErrUserNotFound
		}
		return types.User{}, err
	}

	return u, nil
}

func (s *PostgresUserStore) UserByEmail(email string) (types.User, error) {
	stmt := `
		select p.*, u.id as auth_id, u.email, u.created_at, u.last_sign_in_at
		from public.profiles p
		join auth.users u on u.id = p.id
		where u.email = $1 limit 1;`

	var u types.User
	if err := s.Get(&u, stmt, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.User{}, ErrUserNotFound
		}
		return types.User{}, err
	}

	return u, nil
}
