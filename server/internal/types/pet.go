package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type PetType string

const (
	PetTypeDog     PetType = "Dog"
	PetTypeCat     PetType = "Cat"
	PetTypeUnknown PetType = "Unspecified"
)

type Pet struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    string     `json:"user_id" db:"user_id"`
	Type      *PetType   `json:"type" db:"type"`
	Name      string     `json:"name" db:"name"`
	Tags      PetTags    `json:"tags" db:"tags"`
	DOB       *time.Time `json:"dob" db:"dob"`
	AvatarURI *string    `json:"avatar" db:"avatar_uri"`
	Blurb     *string    `json:"blurb" db:"blurb"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

type PetTags map[string]string

func (t PetTags) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *PetTags) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan PetTag")
	}
	return json.Unmarshal(bytes, t)
}
