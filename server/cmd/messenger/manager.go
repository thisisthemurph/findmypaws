package main

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"paws/internal/repository"
	"paws/internal/types"
	"sync"
)

const recentMessageFetchLimit = 10

var ErrUnauthorized = errors.New("unauthorized")

type Manager struct {
	rooms map[string]*Room
	sync.RWMutex

	//handlers         map[EventType]EventHandler
	conversationRepo repository.ConversationRepository
	logger           *slog.Logger
}

func NewManager(db *sqlx.DB, logger *slog.Logger) *Manager {
	m := &Manager{
		rooms: make(RoomList),
		//handlers:         make(map[EventType]EventHandler),
		conversationRepo: repository.NewConversationsRepository(db),
		logger:           logger,
	}
	//m.setUpHandlers()
	return m
}

//func (m *Manager) HandleEvent(e Event, c *Client) error {
//	handler, ok := m.handlers[e.Type]
//	if !ok {
//		return ErrUnsupportedEventType
//	}
//	return handler(e, c)
//}

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

func (m *Manager) LoadRecentMessages(conversationID int64) ([]types.Message, error) {
	return m.conversationRepo.RecentMessages(conversationID, recentMessageFetchLimit)
}
