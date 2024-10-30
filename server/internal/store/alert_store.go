package store

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"paws/internal/types"
)

var ErrAlertAlreadyExists = errors.New("alert already exists")

type PostgresAlertStore struct {
	*sqlx.DB
}

func NewPostgresAlertStore(db *sqlx.DB) *PostgresAlertStore {
	return &PostgresAlertStore{db}
}

func (s PostgresAlertStore) Create(alert types.Alert) error {
	if exists, err := s.alertExists(alert); err != nil || exists {
		if exists {
			return ErrAlertAlreadyExists
		}
		return err
	}

	stmt := `
		insert into alerts (user_id, anonymous_user_id, pet_id)
		values ($1, $2, $3);`

	var userId, anonymousUserId *string
	if alert.UserId != "" {
		userId = &alert.UserId
	}
	if alert.AnonymousUserId != "" {
		anonymousUserId = &alert.AnonymousUserId
	}

	if _, err := s.Exec(stmt, userId, anonymousUserId, alert.PetId); err != nil {
		return err
	}

	return nil
}

func (s PostgresAlertStore) alertExists(alert types.Alert) (bool, error) {
	stmt := `
		select exists (
			select 1 from alerts
			where user_id = $1
			   or anonymous_user_id = $1
		);`

	var exists bool
	if err := s.Get(&exists, stmt, alert.GetUserId()); err != nil {
		return false, err
	}
	return exists, nil
}
