package response

import (
	"github.com/clerk/clerk-sdk-go/v2"
	"paws/internal/database/model"
	"time"
)

type User struct {
	clerk.User
}

func NewUserFromModel(m clerk.User) User {
	return User{
		User: m,
	}
}

// PrimaryEmailAddress returns the email address associated with the User PrimaryEmailAddressID
// or the first email address, if one is present, if there is no PrimaryEmailAddressID.
func (u User) PrimaryEmailAddress() string {
	if len(u.EmailAddresses) == 0 {
		return ""
	}
	for _, address := range u.EmailAddresses {
		if address.ID == *u.PrimaryEmailAddressID {
			return address.EmailAddress
		}
	}
	return u.EmailAddresses[0].EmailAddress
}

type AnonymousUser struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewAnonymousUserFromModel(m model.AnonymousUser) AnonymousUser {
	return AnonymousUser{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
