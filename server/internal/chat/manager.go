package chat

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"paws/internal/repository"
	"sync"
)

var ErrUnauthorized = errors.New("unauthorized")

// Manager is responsible for managing all the conversation rooms.
type Manager struct {
	rooms map[string]*Room
	sync.RWMutex

	conversationRepo repository.ConversationRepository
	logger           *slog.Logger
}

// NewManager creates an instance of a new manager.
func NewManager(db *sqlx.DB, logger *slog.Logger) *Manager {
	return &Manager{
		rooms:            make(RoomList),
		conversationRepo: repository.NewConversationsRepository(db),
		logger:           logger,
	}
}

// GetOrCreateRoom gets the room if it already exists.
// The conversation is added to the database if it does not exist.
func (m *Manager) GetOrCreateRoom(identifier uuid.UUID, participantID string) (*Room, error) {
	m.Lock()
	defer m.Unlock()

	conversation, err := m.conversationRepo.GetOrCreate(identifier, participantID)
	if err != nil {
		return nil, err
	}

	// Validate that the joining participant is a member of the Room
	if participantID != conversation.PrimaryParticipantID && participantID != conversation.SecondaryParticipantID {
		m.logger.Error("participant not member", "joiningParticipantID", participantID, "conversation", conversation)
		return nil, ErrUnauthorized
	}

	roomKey := NewRoomKey(conversation.ID, conversation.Identifier)
	if r, ok := m.rooms[roomKey.String()]; ok {
		return r, nil
	}

	r := NewRoom(conversation.ID, conversation.Identifier, m)
	m.rooms[roomKey.String()] = r
	go r.run()
	return r, nil
}
