package types

import "time"

type NotificationType string

const SpottedPetNotification NotificationType = "spottedPet"

type Notification struct {
	ID        string            `json:"id"`
	Type      NotificationType  `json:"type"`
	Message   string            `json:"message"`
	Links     map[string]string `json:"links"`
	Seen      bool              `json:"seen"`
	CreatedAt time.Time         `json:"created_at"`
}
