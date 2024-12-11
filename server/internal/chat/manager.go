package chat

import (
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

var ErrUnauthorized = errors.New("unauthorized")

type RoomIdentifier interface {
	ID() int64
	Identifier() uuid.UUID
}

type RoomParticipant interface {
	PrimaryParticipantID() string
	SecondaryParticipantID() string
}

type RoomDetail interface {
	RoomIdentifier
	RoomParticipant
}

type MessageIdentifier interface {
	ID() int64
}

type MessageDetail interface {
	MessageIdentifier
	Text() string
	SenderID() string
	EmojiReaction() string
	CreatedAt() time.Time
}

type ManagerCallbacks struct {
	// HandleRoomCreation is a callback triggered when a room needs to be created or retrieved.
	// The callback function should return the room if it exists, otherwise it should create the room and return it.
	//
	// Returns:
	//   - A RoomDetail instance containing information about the room.
	//   - An error if the room could not be retrieved or created.
	HandleRoomCreation func(identifier uuid.UUID, secondaryParticipantID string) (RoomDetail, error)
	// HandleNewMessage is a callback invoked when a new message is sent in a conversation.
	// If you are persisting messages in a database, you should persist the message in this function and return the message ID.
	//
	// Parameters:
	//   - conversationID: The ID of the conversation containing the message.
	//   - message: The newly created message.
	//
	// Returns:
	//   - The ID of the newly created message.
	//   - An error if the message could not be created.
	HandleNewMessage func(conversationID int64, message NewMessageEvent) (int64, error)
	// HandleEmojiUpdate is a callback allowing you to update the emoji reaction for a specific message.
	//
	// Parameters:
	//   - conversationID: The ID of the conversation containing the message.
	//   - messageID: The ID of the message to update.
	//   - emojiKey: A pointer to a string representing the new emoji reaction. If nil, the emoji should be cleared.
	//
	// Returns:
	//   - An error if the message could not be updated with the emoji key.
	HandleEmojiUpdate func(conversationID, messageID int64, emojiKey *string) error
	// FetchHistoricalMessages is a callback that retrieves historical messages for a given conversation.
	//
	// Parameters:
	//   - conversationID: The unique identifier of the conversation.
	//
	// Returns:
	//   - A slice of MessageDetail instances representing the historical messages.
	//   - An error if the messages could not be retrieved.
	FetchHistoricalMessages func(conversationID int64) ([]MessageDetail, error)
}

// Manager is responsible for managing all the conversation rooms.
type Manager struct {
	rooms map[string]*Room
	sync.RWMutex

	callbacks ManagerCallbacks
	logger    *slog.Logger
}

type ManagerConfig struct {
	Callbacks ManagerCallbacks
	Logger    *slog.Logger
}

// NewManager creates an instance of a new manager.
func NewManager(config ManagerConfig) *Manager {
	return &Manager{
		rooms:     make(RoomList),
		callbacks: config.Callbacks,
		logger:    config.Logger,
	}
}

// GetOrCreateRoom gets the room if it already exists.
// The conversation is added to the database if it does not exist.
func (m *Manager) GetOrCreateRoom(identifier uuid.UUID, participantID string) (*Room, error) {
	m.Lock()
	defer m.Unlock()

	conversation, err := m.callbacks.HandleRoomCreation(identifier, participantID)
	if err != nil {
		return nil, err
	}

	// Validate that the joining participant is a member of the Room
	if participantID != conversation.PrimaryParticipantID() && participantID != conversation.SecondaryParticipantID() {
		m.logger.Error("participant not member", "joiningParticipantID", participantID, "conversation", conversation)
		return nil, ErrUnauthorized
	}

	roomKey := NewRoomKey(conversation.ID(), conversation.Identifier())
	if r, ok := m.rooms[roomKey.String()]; ok {
		return r, nil
	}

	r := NewRoom(conversation.ID(), conversation.Identifier(), m)
	m.rooms[roomKey.String()] = r
	go r.run()
	return r, nil
}
