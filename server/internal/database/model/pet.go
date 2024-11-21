package model

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Pet struct {
	ID        uuid.UUID       `db:"id"`
	UserID    string          `db:"user_id"`
	Type      *string         `db:"type"`
	Name      string          `db:"name"`
	Tags      json.RawMessage `db:"tags"`
	DOB       *time.Time      `db:"dob"`
	AvatarURI *string         `db:"avatar_uri"`
	Blurb     *string         `db:"blurb"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}
