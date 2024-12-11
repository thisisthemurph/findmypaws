package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type EventType string

const (
	EventTypeEmojiReact    EventType = "emoji_react"
	EventTypeNewEmojiReact EventType = "new_emoji_react"
	EventTypeSendMessage   EventType = "send_message"
	EventTypeNewMessage    EventType = "new_message"
	EventTypeTyping        EventType = "typing"
)

var emojiKeyLookup = map[string]string{
	"thumbs-up":     "👍",
	"thumbs-down":   "👎",
	"smiling-face":  "😊",
	"laughing-face": "😆",
	"crying-face":   "😭",
}

var ErrUnsupportedEventType = errors.New("unsupported event type")

type Event struct {
	Type    EventType       `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type SendMessageEvent struct {
	Text     string `json:"text"`
	SenderID string `json:"senderId"`
}

type NewMessageEvent struct {
	SendMessageEvent
	ID            int64     `json:"id"`
	Timestamp     time.Time `json:"timestamp"`
	EmojiReaction *string   `json:"emoji"`
}

type EmojiReactEvent struct {
	EmojiKey       string `json:"emojiKey"`
	ConversationID int64  `json:"conversationId"`
	MessageID      int64  `json:"messageId"`
}

type NewEmojiReactEvent struct {
	MessageID int64   `json:"messageId"`
	Emoji     *string `json:"emoji"`
}

type EventHandler func(e Event, c *Client) error

type eventHandlers struct {
	room   *Room
	logger *slog.Logger
}

func newEventHandlers(room *Room) *eventHandlers {
	return &eventHandlers{
		room:   room,
		logger: room.logger,
	}
}

// SendMessageHandler handles the sending of a new message by a client within a room.
//   - the message is persisted in the database.
//   - an event is sent to each of the room's clients.
func (h *eventHandlers) SendMessageHandler(e Event, c *Client) error {
	h.logger = h.logger.With("handler", "SendMessageHandler")

	var msgEvent SendMessageEvent
	if err := json.Unmarshal(e.Payload, &msgEvent); err != nil {
		return fmt.Errorf("bad payload for %v event: %w", EventTypeSendMessage, err)
	}

	var broadcast NewMessageEvent
	broadcast.SendMessageEvent = msgEvent
	broadcast.Timestamp = time.Now()

	messageID, err := h.room.manager.callbacks.HandleNewMessage(h.room.key.ConversationID, broadcast)
	if err != nil {
		h.logger.Error("error persisting message in database", "error", err)
	}

	broadcast.ID = messageID
	data, err := json.Marshal(broadcast)
	if err != nil {
		return fmt.Errorf("could not marshal new message: %w", err)
	}

	var outgoingEvent Event
	outgoingEvent.Type = EventTypeNewMessage
	outgoingEvent.Payload = data

	for client := range c.room.clients {
		client.egress <- outgoingEvent
	}
	return nil
}

// EmojiReactHandler handles emoji reactions to messages.
//   - the database is updated to reflect the new emoji reaction.
//   - an event is sent to all room clients to notify them of the emoji reaction.
func (h *eventHandlers) EmojiReactHandler(e Event, c *Client) error {
	h.logger = h.logger.With("handler", "EmojiReactHandler")

	var emojiEvent EmojiReactEvent
	if err := json.Unmarshal(e.Payload, &emojiEvent); err != nil {
		return fmt.Errorf("bad payload for %v event: %w", EventTypeEmojiReact, err)
	}
	h.logger.Debug("emoji", "event", emojiEvent)

	var emojiKey *string
	if emojiEvent.EmojiKey == "" {
		emojiKey = nil
	} else {
		emojiKey = &emojiEvent.EmojiKey
	}

	if err := h.room.manager.callbacks.HandleEmojiUpdate(emojiEvent.ConversationID, emojiEvent.MessageID, emojiKey); err != nil {
		return fmt.Errorf("could not update message emoji: %w", err)
	}

	var emoji *string
	if selectedEmoji, ok := emojiKeyLookup[emojiEvent.EmojiKey]; ok {
		emoji = &selectedEmoji
	}

	var outgoingEvent Event
	outgoingEvent.Type = EventTypeNewEmojiReact
	eventData := NewEmojiReactEvent{
		MessageID: emojiEvent.MessageID,
		Emoji:     emoji,
	}
	data, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("could not marshal new new emji reaction event: %w", err)
	}

	outgoingEvent.Payload = data
	for client := range c.room.clients {
		client.egress <- outgoingEvent
	}
	return nil
}
