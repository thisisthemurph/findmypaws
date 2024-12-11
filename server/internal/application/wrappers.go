package application

import (
	"time"

	"github.com/google/uuid"
	"paws/internal/database/model"
)

type ConversationWrapper struct {
	*model.Conversation
}

func (cw ConversationWrapper) ID() int64 {
	return cw.Conversation.ID
}

func (cw ConversationWrapper) Identifier() uuid.UUID {
	return cw.Conversation.Identifier
}

func (cw ConversationWrapper) PrimaryParticipantID() string {
	return cw.Conversation.PrimaryParticipantID
}

func (cw ConversationWrapper) SecondaryParticipantID() string {
	return cw.Conversation.SecondaryParticipantID
}

type MessageWrapper struct {
	*model.Message
}

func (mw MessageWrapper) ID() int64 {
	return mw.Message.ID
}

func (mw MessageWrapper) Text() string {
	return mw.Message.Text
}

func (mw MessageWrapper) SenderID() string {
	return mw.Message.SenderID
}

func (mw MessageWrapper) EmojiReaction() string {
	if mw.Message.EmojiReaction == nil {
		return ""
	}
	return *mw.Message.EmojiReaction
}

func (mw MessageWrapper) CreatedAt() time.Time {
	return mw.Message.CreatedAt
}
