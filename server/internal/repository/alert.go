package repository

import (
	"github.com/jmoiron/sqlx"
	"paws/internal/types"
)

type AlertRepository interface {
	List(userID string) ([]types.AlertModel, error)
	Create(alert types.AlertIdentifiers) error
	MarkAllAsRead(userID string) error
}

type PostgresAlertRepository struct {
	db *sqlx.DB
}

func NewPostgresAlertRepository(db *sqlx.DB) *PostgresAlertRepository {
	return &PostgresAlertRepository{db}
}

func (r *PostgresAlertRepository) List(userID string) ([]types.AlertModel, error) {
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
	if err := r.db.Select(&aa, stmt, userID); err != nil {
		return nil, err
	}
	return aa, nil
}

func (r *PostgresAlertRepository) Create(alert types.AlertIdentifiers) error {
	if exists, err := r.alertExists(alert); err != nil || exists {
		if exists {
			return ErrAlreadyExists
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

	if _, err := r.db.Exec(stmt, userId, anonymousUserId, alert.PetId); err != nil {
		return err
	}

	return nil
}

func (r *PostgresAlertRepository) MarkAllAsRead(userID string) error {
	stmt := `
		update alerts set seen_at = NOW()
		from pets
		where alerts.pet_id = pets.id
		and pets.user_id = $1;`

	_, err := r.db.Exec(stmt, userID)
	return err
}

func (r *PostgresAlertRepository) alertExists(alert types.Alerter) (bool, error) {
	stmt := `
		select exists (
			select 1 from alerts
			where pet_id = $1 and (
				user_id = $2
			   	or anonymous_user_id = $2
			)
		);`

	var exists bool
	if err := r.db.Get(&exists, stmt, alert.GetPetId(), alert.GetReporterUserId()); err != nil {
		return false, err
	}
	return exists, nil
}
