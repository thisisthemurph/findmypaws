package types

import (
	"github.com/google/uuid"
	"time"
)

type Alerter interface {
	// IsAnonymous returns true if the reporter is anonymous.
	IsAnonymous() bool
	// GetReporterUserId returns the ID of the reporter, anonymous or user.
	GetReporterUserId() string
	GetPetId() uuid.UUID
}

type AlertIdentifiers struct {
	PetId           uuid.UUID `db:"pet_id"`
	UserID          string    `db:"user_id"`
	AnonymousUserId string    `db:"anonymous_user_id"`
}

type AlertModel struct {
	ID        int64      `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	SeenAt    *time.Time `db:"seen_at"`

	AlertIdentifiers
}

func (a AlertIdentifiers) IsAnonymous() bool {
	return a.UserID == ""
}

func (a AlertIdentifiers) GetReporterUserId() string {
	if a.IsAnonymous() {
		return a.AnonymousUserId
	}
	return a.UserID
}

func (a AlertIdentifiers) GetPetId() uuid.UUID {
	return a.PetId
}
