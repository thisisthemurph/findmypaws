package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type EventType string

const (
	EventTypeSendMessage EventType = "send_message"
	EventTypeNewMessage  EventType = "new_message"
	EventTypeTyping      EventType = "typing"
)

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
	Timestamp time.Time `json:"timestamp"`
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

	if err := r.PersistMessage(broadcast); err != nil {
		r.logger.Error("error persisting message in database", "error", err)
	}
	return nil
}
