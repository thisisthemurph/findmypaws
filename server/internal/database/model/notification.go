package model

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Notification struct {
	ID        int64           `db:"id"`
	UserID    string          `db:"user_id"`
	PetID     uuid.UUID       `db:"pet_id"`
	Type      string          `db:"type"`
	Detail    json.RawMessage `db:"detail"`
	CreatedAt time.Time       `db:"created_at"`
	SeenAt    *time.Time      `db:"seen_at"`
}

type SpottedPetNotificationDetail struct {
	SpotterName string    `json:"spotter_name"`
	IsAnonymous bool      `json:"is_anonymous"`
	PetName     string    `json:"pet_name"`
	PetID       uuid.UUID `json:"pet_id"`
}

func NewSpottedPetNotification(userID string, detail SpottedPetNotificationDetail) (Notification, error) {
	if detail.PetID == uuid.Nil {
		return Notification{}, fmt.Errorf("invalid spotted pet notification detail; PetId is required")
	}

	detailJSON, err := json.Marshal(detail)
	if err != nil {
		return Notification{}, fmt.Errorf("error marshalling notification detail: %w", err)
	}

	return Notification{
		UserID:    userID,
		PetID:     detail.PetID,
		Type:      "spotted_pet",
		Detail:    detailJSON,
		CreatedAt: time.Now(),
	}, nil
}
