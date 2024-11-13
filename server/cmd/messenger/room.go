package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"paws/internal/types"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 1024
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     checkOrigin,
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

type RoomKey struct {
	ConversationID int64
	Identifier     uuid.UUID
}

func (rk RoomKey) String() string {
	return fmt.Sprintf("room:%v:%d", rk.Identifier, rk.ConversationID)
}

type RoomList map[string]*Room

func NewRoomKey(conversationID int64, identifier uuid.UUID) RoomKey {
	return RoomKey{
		ConversationID: conversationID,
		Identifier:     identifier,
	}
}

type Room struct {
	key     RoomKey
	manager *Manager
	logger  *slog.Logger

	handlers map[EventType]EventHandler
	clients  ClientList
	sync.RWMutex

	join    chan *Client
	leave   chan *Client
	forward chan Event
}

func NewRoom(conversationID int64, identifier uuid.UUID, manager *Manager) *Room {
	roomKey := NewRoomKey(conversationID, identifier)
	room := &Room{
		key:     roomKey,
		manager: manager,
		logger:  manager.logger.With("roomKey", roomKey.String()),

		handlers: make(map[EventType]EventHandler),
		clients:  make(map[*Client]struct{}),
		join:     make(chan *Client),
		leave:    make(chan *Client),
		forward:  make(chan Event),
	}
	room.setUpHandlers()
	return room
}

func (r *Room) ServeWS(w http.ResponseWriter, req *http.Request) {
	roomID := req.URL.Query().Get("r")
	if roomID == "" {
		http.Error(w, "Room key required", http.StatusBadRequest)
		return
	}

	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		r.logger.Error("error upgrading socket", "error", err)
		return
	}

	client := NewClient(socket, r)

	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}

func (r *Room) HandleEvent(e Event, c *Client) error {
	handler, ok := r.handlers[e.Type]
	if !ok {
		return ErrUnsupportedEventType
	}
	return handler(e, c)
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.join:
			r.logger.Info("join", "Client", client)
			r.addClient(client)
			if err := r.EgressHistoricalMessages(client); err != nil {
				r.logger.Error("failed to egress historical messages", "error", err)
			}
		case client := <-r.leave:
			r.logger.Info("leave", "Client", client)
			r.removeClient(client)
		case message := <-r.forward:
			r.logger.Info("forward", "roomID", r.key, "msg", message)
			for client := range r.clients {
				client.egress <- message
			}
		}
	}
}

func (r *Room) addClient(client *Client) {
	r.Lock()
	defer r.Unlock()
	r.clients[client] = struct{}{}
}

func (r *Room) removeClient(client *Client) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.clients[client]; ok {
		close(client.egress)
		if err := client.socket.Close(); err != nil {
			r.logger.Error("error closing client socket", "client", client, "error", err)
		}
		delete(r.clients, client)
	}
}

func (r *Room) PersistMessage(message NewMessageEvent) error {
	m := &types.Message{
		ConversationID: r.key.ConversationID,
		SenderID:       message.SenderID,
		Text:           message.Text,
	}

	return r.manager.conversationRepo.CreateMessage(m)
}

func (r *Room) EgressHistoricalMessages(client *Client) error {
	messages, err := r.manager.conversationRepo.ListHistoricalMessages(r.key.ConversationID, time.Now(), 10)
	if err != nil {
		return fmt.Errorf("error querying historical messages: %w", err)
	}

	for _, message := range messages {
		msg := NewMessageEvent{}
		msg.Text = message.Text
		msg.SenderID = message.SenderID
		msg.Timestamp = message.CreatedAt
		messageJSON, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("error marshalling message: %w", err)
		}

		client.egress <- Event{
			Type:    EventTypeNewMessage,
			Payload: messageJSON,
		}
	}

	return nil
}

func (r *Room) setUpHandlers() {
	r.handlers[EventTypeSendMessage] = r.SendMessageHandler
}

func checkOrigin(r *http.Request) bool {
	return true
}
