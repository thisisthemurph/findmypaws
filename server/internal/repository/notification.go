package repository

import (
	"encoding/json"
	"fmt"
	"paws/internal/database/model"
	"paws/internal/response"

	"github.com/jmoiron/sqlx"
)

type NotificationRepository interface {
	List(userID string) ([]model.Notification, error)
	Create(n *model.Notification) error
	MarkAllSeen(userID string) error
	RecentlyNotified(n model.Notification) (bool, error)
}

type postgresNotificationRepository struct {
	db *sqlx.DB
}

func NewNotificationRepository(db *sqlx.DB) NotificationRepository {
	return &postgresNotificationRepository{
		db: db,
	}
}

func (r *postgresNotificationRepository) List(userID string) ([]model.Notification, error) {
	q := `
		select * from notifications 
		where user_id = $1
		order by created_at desc;`

	var nn []model.Notification
	if err := r.db.Select(&nn, q, userID); err != nil {
		return nil, err
	}
	return nn, nil
}

func (r *postgresNotificationRepository) Create(n *model.Notification) error {
	stmt := `
		insert into notifications (user_id, type, detail)
		values ($1, $2, $3)
		returning id, created_at;`

	if err := r.db.Get(n, stmt, n.UserID, n.Type, n.Detail); err != nil {
		return err
	}
	return nil
}

func (r *postgresNotificationRepository) MarkAllSeen(userID string) error {
	_, err := r.db.Exec("update notifications set seen_at = now() where user_id = $1;", userID)
	return err
}

func (r *postgresNotificationRepository) RecentlyNotified(n model.Notification) (bool, error) {
	switch n.Type {
	case "spotted_pet":
		return r.spottedPetNotificationRecentlyNotified(n)
	default:
		return false, fmt.Errorf("unknown notification type %v", n.Type)
	}
}

func (r *postgresNotificationRepository) spottedPetNotificationRecentlyNotified(n model.Notification) (bool, error) {
	var detail response.SpottedPetNotificationDetail
	if err := json.Unmarshal(n.Detail, &detail); err != nil {
		return false, err
	}

	q := `
		select exists(
			select 1 from notifications
			where type = 'spotted_pet'
			  and created_at >= now() - interval '1 day'
			  and detail ->> 'pet_id' = $1
		);`

	var exists bool
	if err := r.db.QueryRow(q, detail.PetID).Scan(&exists); err != nil {
		return exists, err
	}
	return exists, nil
}
