package types

import (
	"github.com/clerk/clerk-sdk-go/v2"
	"time"
)

type User struct {
	clerk.User
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
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
