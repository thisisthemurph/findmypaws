package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log/slog"
	"time"
)

const pongWait = 10 * time.Second

type ClientList map[*Client]struct{}

type Client struct {
	room   *Room
	socket *websocket.Conn
	egress chan Event
	logger *slog.Logger
}

func NewClient(ws *websocket.Conn, room *Room) *Client {
	return &Client{
		room:   room,
		socket: ws,
		egress: make(chan Event, messageBufferSize),
		logger: room.logger,
	}
}

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
			c.logger.Error("error unmarshalling message", "error", err)
			break
		}

		if err := c.room.HandleEvent(ev, c); err != nil {
			c.logger.Error("error handling event", "type", ev.Type, "error", err)
		}
	}
}

func (c *Client) write() {
	for m := range c.egress {
		data, err := json.Marshal(m)
		if err != nil {
			c.logger.Error("error parsing m JSON", "error", err)
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
