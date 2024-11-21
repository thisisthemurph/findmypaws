package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/jmoiron/sqlx"
	"paws/internal/database/model"
)

type UserRepository interface {
	GetUser(id string) (clerk.User, error)
	UpsertUser(u clerk.User) error
	DeleteUser(id string) error
	GetAnonymousUser(id string) (model.AnonymousUser, error)
	UpsertAnonymousUser(u *model.AnonymousUser) error
}

type postgresUserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &postgresUserRepository{
		db: db,
	}
}

func (r *postgresUserRepository) GetUser(id string) (clerk.User, error) {
	stmt := "select data from users where id = $1;"
	var b []byte
	if err := r.db.QueryRow(stmt, id).Scan(&b); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return clerk.User{}, ErrNotFound
		}
		return clerk.User{}, fmt.Errorf("error getting user from database: %w", err)
	}

	var usr clerk.User
	if err := json.Unmarshal(b, &usr); err != nil {
		return clerk.User{}, fmt.Errorf("error unmarshalling user JSON: %w", err)
	}
	return usr, nil
}

func (r *postgresUserRepository) UpsertUser(u clerk.User) error {
	stmt := `
		insert into users (id, data) 
		values ($1, $2)
		on conflict (id) do update
			set data = $2;`

	data, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("error marshalling user JSON: %w", err)
	}

	if _, err := r.db.Exec(stmt, u.ID, data); err != nil {
		return fmt.Errorf("error upserting user: %w", err)
	}
	return nil
}

func (r *postgresUserRepository) DeleteUser(id string) error {
	stmt := "delete from users where id = $1;"
	if _, err := r.db.Exec(stmt, id); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}

func (r *postgresUserRepository) GetAnonymousUser(id string) (model.AnonymousUser, error) {
	var u model.AnonymousUser
	if err := r.db.Get(&u, "SELECT * FROM anonymous_users WHERE id = $1", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.AnonymousUser{}, ErrNotFound
		}
		return model.AnonymousUser{}, err
	}
	return u, nil
}

func (r *postgresUserRepository) UpsertAnonymousUser(u *model.AnonymousUser) error {
	q := `
		insert into anonymous_users (id, name)
		values ($1, $2)
		on conflict (id) do update set name = $2
		returning created_at, updated_at;`

	err := r.db.Get(u, q, u.ID, u.Name)
	return err
}
