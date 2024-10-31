package types

import "github.com/google/uuid"

type Alert struct {
	AnonymousUserId string
	UserId          string
	PetId           uuid.UUID
}

func (a Alert) IsAnonymous() bool {
	return a.AnonymousUserId != ""
}

func (a Alert) GetUserId() string {
	if a.IsAnonymous() {
		return a.AnonymousUserId
	}
	return a.UserId
}
