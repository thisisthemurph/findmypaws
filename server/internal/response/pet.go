package response

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"paws/internal/database/model"
	"time"

	"github.com/google/uuid"
)

type PetType string

const (
	PetTypeDog     PetType = "Dog"
	PetTypeCat     PetType = "Cat"
	PetTypeUnknown PetType = "Unspecified"
)

func NewPetType(v *string) PetType {
	switch *v {
	case "Dog":
		return PetTypeDog
	case "Cat":
		return PetTypeCat
	default:
		return PetTypeUnknown
	}
}

type Pet struct {
	ID        uuid.UUID  `json:"id"`
	UserID    string     `json:"user_id"`
	Type      PetType    `json:"type"`
	Name      string     `json:"name"`
	Tags      PetTags    `json:"tags"`
	DOB       *time.Time `json:"dob"`
	AvatarURI *string    `json:"avatar"`
	Blurb     *string    `json:"blurb"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewPetFromModel(m *model.Pet) Pet {
	var tags PetTags
	if err := json.Unmarshal(m.Tags, &tags); err != nil {
		tags = make(PetTags)
	}

	p := Pet{
		ID:        m.ID,
		UserID:    m.UserID,
		Type:      NewPetType(m.Type),
		Name:      m.Name,
		Tags:      tags,
		DOB:       m.DOB,
		AvatarURI: m.AvatarURI,
		Blurb:     m.Blurb,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	return p
}

type PetTags map[string]string

func NewPetTags(j json.RawMessage) PetTags {
	var tags PetTags
	if err := json.Unmarshal(j, &tags); err != nil {
		tags = make(PetTags)
	}
	return tags
}

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
