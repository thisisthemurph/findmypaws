package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrNotAuthorized = errors.New("not authorized")
)

type Repositories struct {
	PetRepository          PetRepository
	NotificationRepository NotificationRepository
	ConversationRepository ConversationRepository
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		PetRepository:          NewPetRepository(db),
		NotificationRepository: NewNotificationRepository(db),
		ConversationRepository: NewConversationsRepository(db),
	}
}
