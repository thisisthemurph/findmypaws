package chat

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

const pongWait = 10 * time.Second

// ClientList represents a list of Client.
type ClientList map[*Client]struct{}

// Client represents a user within a Room.
// A Room can have only one instance of each Client
// A user can be in multiple rooms, represented by different clients.
type Client struct {
	room   *Room
	socket *websocket.Conn
	egress chan Event
	logger *slog.Logger
}

// NewClient creates an instance of a Client.
func NewClient(ws *websocket.Conn, room *Room) *Client {
	return &Client{
		room:   room,
		socket: ws,
		egress: make(chan Event, messageBufferSize),
		logger: room.logger,
	}
}

// read starts an infinite loop for the client, checking for new messages on the client's socket.
// If a message is received, it is unmarshalled and sent to the room for handling.
func (c *Client) read() {
	defer func() {
		c.room.removeClient(c)
	}()

	c.socket.SetReadLimit(messageBufferSize)
	//if err := c.socket.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
	//	c.logger.Error("read deadline error", "error", err)
	//	return
	//}
	//c.socket.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.socket.ReadMessage()
		if err != nil {
			// Handle bad closed connections
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error("error reading message", "error", err)
			}
			break
		}

		var ev Event
		if err := json.Unmarshal(payload, &ev); err != nil {
			c.logger.Error("error unmarshalling event payload", "payload", payload, "error", err)
			continue
		}

		if err := c.room.HandleEvent(ev, c); err != nil {
			c.logger.Error("error handling event", "type", ev.Type, "error", err)
		}
	}
}

func (c *Client) write() {
	for event := range c.egress {
		data, err := json.Marshal(event)
		if err != nil {
			c.logger.Error("error parsing event JSON", "error", err)
			return
		}
		if err := c.socket.WriteMessage(websocket.TextMessage, data); err != nil {
			return
		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	slog.Debug("pong")
	return c.socket.SetReadDeadline(time.Now().Add(pongWait))
}
