package main

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (r *Room) SendMessageHandler(e Event, c *Client) error {
	var msgEvent SendMessageEvent
	if err := json.Unmarshal(e.Payload, &msgEvent); err != nil {
		return fmt.Errorf("bad payload for %v event: %w", EventTypeSendMessage, err)
	}

	var broadcast NewMessageEvent
	broadcast.SendMessageEvent = msgEvent
	broadcast.Timestamp = time.Now()

	message, err := r.PersistMessage(broadcast)
	if err != nil {
		r.logger.Error("error persisting message in database", "error", err)
	}

	broadcast.ID = message.ID
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

func (r *Room) EmojiReactHandler(e Event, c *Client) error {
	var emojiEvent EmojiReactEvent
	if err := json.Unmarshal(e.Payload, &emojiEvent); err != nil {
		return fmt.Errorf("bad payload for %v event: %w", EventTypeEmojiReact, err)
	}
	r.logger.Info("emoji", "event", emojiEvent)

	message, err := r.manager.conversationRepo.GetMessage(emojiEvent.ConversationID, emojiEvent.MessageID)
	if err != nil {
		return fmt.Errorf("could not get message from conversation: %w", err)
	}

	if emojiEvent.EmojiKey == "" {
		message.EmojiReaction = nil
	} else {
		message.EmojiReaction = &emojiEvent.EmojiKey
	}

	if err := r.manager.conversationRepo.UpdateMessage(message); err != nil {
		return fmt.Errorf("could not update message: %w", err)
	}

	var emoji *string
	selectedEmoji, ok := emojiKeyLookup[emojiEvent.EmojiKey]
	if ok {
		emoji = &selectedEmoji
	}

	var outgoingEvent Event
	outgoingEvent.Type = EventTypeNewEmojiReact
	eventData := NewEmojiReactEvent{
		MessageID: message.ID,
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
