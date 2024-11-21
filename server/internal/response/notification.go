package response

import (
	"encoding/json"
	"fmt"
	"paws/internal/database/model"
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	SpottedPetNotification NotificationType = "spotted_pet"
	UnknownNotification    NotificationType = "unknown"
)

func NewNotificationType(t string) NotificationType {
	switch t {
	case "spotted_pet":
		return SpottedPetNotification
	default:
		return UnknownNotification
	}
}

// Notification represents a generic user facing notification.
type Notification struct {
	ID        string           `json:"id"`
	Type      NotificationType `json:"type"`
	Message   string           `json:"message"`
	Link      string           `json:"link"`
	CreatedAt time.Time        `json:"created_at"`
	Seen      bool             `json:"seen"`
}

func NewNotificationFromModel(m model.Notification) (Notification, bool) {
	detail, err := parseNotificationDetail(m)
	if err != nil {
		return Notification{}, false
	}

	notification := Notification{
		ID:        fmt.Sprintf("%s_%d", m.Type, m.ID),
		Type:      NewNotificationType(m.Type),
		Message:   detail.Message(),
		Link:      detail.Link(),
		CreatedAt: m.CreatedAt,
		Seen:      m.SeenAt != nil,
	}
	return notification, true
}

func parseNotificationDetail(m model.Notification) (NotificationDetail, error) {
	var detail NotificationDetail

	switch m.Type {
	case "spotted_pet":
		var messageDetail SpottedPetNotificationDetail
		if err := json.Unmarshal(m.Detail, &messageDetail); err != nil {
			return nil, err
		}
		detail = messageDetail
	default:
		return nil, fmt.Errorf("unknown notification type: %s", m.Type)
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
