package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	SpottedPetNotification NotificationType = "spotted_pet"
)

// Notification represents a generic user facing notification.
type Notification struct {
	ID        string           `json:"id"`
	Type      NotificationType `json:"type"`
	Message   string           `json:"message"`
	Link      string           `json:"link"`
	CreatedAt time.Time        `json:"created_at"`
	Seen      bool             `json:"seen"`
}

// NotificationModel represents a row in the notifications table.
type NotificationModel struct {
	ID        int64            `db:"id"`
	UserID    string           `db:"user_id"`
	PetID     uuid.UUID        `db:"pet_id"`
	Type      NotificationType `db:"type"`
	Detail    json.RawMessage  `db:"detail"`
	CreatedAt time.Time        `db:"created_at"`
	SeenAt    *time.Time       `db:"seen_at"`
}

// Notification converts a NotificationModel to a Notification.
func (n NotificationModel) Notification() (Notification, bool) {
	detail, err := n.parseDetail()
	if err != nil {
		return Notification{}, false
	}

	notification := Notification{
		ID:        fmt.Sprintf("%s_%d", n.Type, n.ID),
		Type:      n.Type,
		Message:   detail.Message(),
		Link:      detail.Link(),
		CreatedAt: n.CreatedAt,
		Seen:      n.SeenAt != nil,
	}
	return notification, true
}

func (n NotificationModel) parseDetail() (NotificationDetail, error) {
	var detail NotificationDetail

	switch n.Type {
	case "spotted_pet":
		var messageDetail SpottedPetNotificationDetail
		if err := json.Unmarshal(n.Detail, &messageDetail); err != nil {
			return nil, err
		}
		detail = messageDetail
	default:
		return nil, fmt.Errorf("unknown notification type: %s", n.Type)
	}

	return detail, nil
}

type NotificationDetail interface {
	Message() string
	Link() string
}

type SpottedPetNotificationDetail struct {
	SpotterName string    `json:"spotter_name"`
	IsAnonymous bool      `json:"is_anonymous"`
	PetName     string    `json:"pet_name"`
	PetID       uuid.UUID `json:"pet_id"`
}

func (d SpottedPetNotificationDetail) Message() string {
	if d.IsAnonymous {
		return fmt.Sprintf("%s has been spotted by an anonymous user", d.PetName)
	}
	return fmt.Sprintf("%s has been spotted by %s", d.PetName, d.SpotterName)
}

func (d SpottedPetNotificationDetail) Link() string {
	return fmt.Sprintf("/pet/%s", d.PetID)
}
