package response

import (
	"github.com/google/uuid"
	"paws/internal/database/model"
	"time"
)

func NewConversationFromModel(m model.Conversation) Conversation {
	return Conversation{
		ID:                     m.ID,
		Identifier:             m.Identifier,
		PrimaryParticipantID:   m.PrimaryParticipantID,
		SecondaryParticipantID: m.SecondaryParticipantID,
		LastMessageAt:          m.LastMessageAt,
		CreatedAt:              m.CreatedAt,
	}
}

type Conversation struct {
	ID                     int64      `json:"id"`
	Identifier             uuid.UUID  `json:"identifier"`
	PrimaryParticipantID   string     `json:"primaryParticipantId"`
	SecondaryParticipantID string     `json:"secondaryParticipantId"`
	LastMessageAt          *time.Time `json:"lastMessageAt"`
	CreatedAt              time.Time  `json:"createdAt"`
}

func NewMessageFromModel(m model.Message) Message {
	return Message{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		SenderID:       m.SenderID,
		Text:           m.Text,
		EmojiReaction:  m.EmojiReaction,
		CreatedAt:      m.CreatedAt,
		ReadAt:         m.ReadAt,
	}
}

type Message struct {
	ID             int64      `json:"id"`
	ConversationID int64      `json:"conversationId"`
	SenderID       string     `json:"senderId"`
	Text           string     `json:"text"`
	EmojiReaction  *string    `json:"emojiReaction"`
	CreatedAt      time.Time  `json:"createdAt"`
	ReadAt         *time.Time `json:"readAt"`
}
