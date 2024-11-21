package model

import (
	"github.com/google/uuid"
	"time"
)

type Conversation struct {
	ID                     int64      `db:"id"`
	Identifier             uuid.UUID  `db:"identifier"`
	PrimaryParticipantID   string     `db:"primary_participant_id"`
	SecondaryParticipantID string     `db:"secondary_participant_id"`
	LastMessageAt          *time.Time `db:"last_message_at"`
	CreatedAt              time.Time  `db:"created_at"`
}

type Message struct {
	ID             int64      `db:"id"`
	ConversationID int64      `db:"conversation_id"`
	SenderID       string     `db:"sender_id"`
	Text           string     `db:"text"`
	EmojiReaction  *string    `db:"emoji_reaction"`
	CreatedAt      time.Time  `db:"created_at"`
	ReadAt         *time.Time `db:"read_at"`
}
