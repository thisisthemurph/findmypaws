package repository

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"paws/internal/types"
)

type UserRepository interface {
	GetAnonymousUser(id string) (types.AnonymousUser, error)
	UpsertAnonymousUser(u *types.AnonymousUser) error
}

type postgresUserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &postgresUserRepository{
		db: db,
	}
}

func (r *postgresUserRepository) GetAnonymousUser(id string) (types.AnonymousUser, error) {
	var u types.AnonymousUser
	if err := r.db.Get(&u, "SELECT * FROM anonymous_users WHERE id = $1", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.AnonymousUser{}, ErrNotFound
		}
		return types.AnonymousUser{}, err
	}
	return u, nil
}

func (r *postgresUserRepository) UpsertAnonymousUser(u *types.AnonymousUser) error {
	q := `
		insert into anonymous_users (id, name)
		values ($1, $2)
		on conflict (id) do update set name = $2
		returning created_at, updated_at;`

	err := r.db.Get(u, q, u.ID, u.Name)
	return err
}
