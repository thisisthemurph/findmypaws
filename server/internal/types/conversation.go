package types

import (
	"github.com/google/uuid"
	"time"
)

type Conversation struct {
	ID                     int64      `db:"id" json:"id"`
	Identifier             uuid.UUID  `db:"identifier" json:"identifier"`
	PrimaryParticipantID   string     `db:"primary_participant_id" json:"primaryParticipantId"`
	SecondaryParticipantID string     `db:"secondary_participant_id" json:"secondaryParticipantId"`
	LastMessageAt          *time.Time `db:"last_message_at" json:"lastMessageAt"`
	CreatedAt              time.Time  `db:"created_at" json:"createdAt"`
}

type Message struct {
	ID             int64      `db:"id" json:"id"`
	ConversationID int64      `db:"conversation_id" json:"conversationId"`
	SenderID       string     `db:"sender_id" json:"senderId"`
	Text           string     `db:"text" json:"text"`
	CreatedAt      time.Time  `db:"created_at" json:"createdAt"`
	ReadAt         *time.Time `db:"read_at" json:"readAt"`
}
