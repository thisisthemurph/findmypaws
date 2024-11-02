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

func (s PostgresAlertStore) Alerts(userID string) ([]types.AlertModel, error) {
	stmt := `
		select
		    a.id, 
		    a.pet_id,
		    coalesce(a.user_id, '') as user_id, 
		    coalesce(a.anonymous_user_id, '') as anonymous_user_id, 
		    a.created_at,
		    a.seen_at
		from alerts a
		join pets p on a.pet_id = p.id
		where p.user_id = $1;`

	var aa []types.AlertModel
	if err := s.Select(&aa, stmt, userID); err != nil {
		return nil, err
	}
	return aa, nil
}

func (s PostgresAlertStore) Create(alert types.AlertIdentifiers) error {
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
	if alert.UserID != "" {
		userId = &alert.UserID
	}
	if alert.AnonymousUserId != "" {
		anonymousUserId = &alert.AnonymousUserId
	}

	if _, err := s.Exec(stmt, userId, anonymousUserId, alert.PetId); err != nil {
		return err
	}

	return nil
}

func (s PostgresAlertStore) MarkAllAsRead(userID string) error {
	stmt := `
		update alerts set seen_at = NOW()
		from pets
		where alerts.pet_id = pets.id
		and pets.user_id = $1;`

	_, err := s.Exec(stmt, userID)
	return err
}

func (s PostgresAlertStore) alertExists(alert types.Alerter) (bool, error) {
	stmt := `
		select exists (
			select 1 from alerts
			where pet_id = $1 and (
				user_id = $2
			   	or anonymous_user_id = $2
			)
		);`

	var exists bool
	if err := s.Get(&exists, stmt, alert.GetPetId(), alert.GetReporterUserId()); err != nil {
		return false, err
	}
	return exists, nil
}
