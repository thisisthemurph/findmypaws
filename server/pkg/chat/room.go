package chat

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

// String returns a string representation of the RoomKey.
func (rk RoomKey) String() string {
	return fmt.Sprintf("room:%v:%d", rk.Identifier, rk.ConversationID)
}

type RoomList map[string]*Room

// NewRoomKey creates a new instance of a RoomKey with the given conversation ID and identifier.
func NewRoomKey(conversationID int64, identifier uuid.UUID) RoomKey {
	return RoomKey{
		ConversationID: conversationID,
		Identifier:     identifier,
	}
}

// Room represents a room for a single conversation.
type Room struct {
	key     RoomKey
	manager *Manager
	logger  *slog.Logger

	clients ClientList
	sync.RWMutex

	join    chan *Client
	leave   chan *Client
	forward chan Event

	handlers *eventHandlers
}

// NewRoom instantiates a new Room.
func NewRoom(conversationID int64, identifier uuid.UUID, manager *Manager) *Room {
	roomKey := NewRoomKey(conversationID, identifier)
	room := &Room{
		key:     roomKey,
		manager: manager,
		logger:  manager.logger.With("roomKey", roomKey.String()),

		clients: make(map[*Client]struct{}),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		forward: make(chan Event),
	}
	handlers := newEventHandlers(room)
	room.handlers = handlers
	return room
}

// ServeWS takes the initial HTTP request and updates it to a WebSocket connection.
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

// HandleEvent performs the appropriate action for the Event depending on the event type.
func (r *Room) HandleEvent(e Event, c *Client) error {
	switch e.Type {
	case EventTypeSendMessage:
		return r.handlers.SendMessageHandler(e, c)
	case EventTypeEmojiReact:
		return r.handlers.EmojiReactHandler(e, c)
	case EventTypeTyping:
		return r.handlers.SendTypingIndication(e, c)
	default:
		return ErrUnsupportedEventType
	}
}

// Run runs the room in an infinite loop handling any events that occur.
func (r *Room) run() {
	for {
		select {
		case client := <-r.join:
			r.logger.Debug("join", "Client", client)
			r.addClient(client)
			if err := r.EgressHistoricalMessages(client); err != nil {
				r.logger.Error("failed to egress historical messages", "error", err)
			}
		case client := <-r.leave:
			r.logger.Debug("leave", "Client", client)
			r.removeClient(client)
		case message := <-r.forward:
			r.logger.Debug("forward", "roomID", r.key, "msg", message)
			for client := range r.clients {
				client.egress <- message
			}
		}
	}
}

// addClient adds a new Client (user) to the room.
// The client will be overridden if they already exist.
func (r *Room) addClient(client *Client) {
	r.Lock()
	defer r.Unlock()
	r.clients[client] = struct{}{}
}

// removeClient removes a client (user) from the room if they exist.
// Nothing happens if the client is not in the room.
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

// EgressHistoricalMessages sends historical messages to a specific client (user). A client belongs to a specific room.
func (r *Room) EgressHistoricalMessages(client *Client) error {
	messages, err := r.manager.callbacks.FetchHistoricalMessages(r.key.ConversationID)
	if err != nil {
		return fmt.Errorf("error querying historical messages: %w", err)
	}

	for _, message := range messages {
		msg := NewMessageEvent{}
		msg.ID = message.ID()
		msg.Text = message.Text()
		msg.SenderID = message.SenderID()
		msg.Timestamp = message.CreatedAt()

		if message.EmojiReaction() != "" {
			emoji, ok := emojiKeyLookup[message.EmojiReaction()]
			if ok {
				msg.EmojiReaction = &emoji
			}
		}

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

// checkOrigin determines if the requesting URL is permitted.
func checkOrigin(r *http.Request) bool {
	return true
}
