package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type WebSocketConnection interface {
	Close() error
	ReadMessage() (messageType int, p []byte, err error)
	SetPongHandler(h func(appData string) error)
	SetReadDeadline(t time.Time) error
	SetReadLimit(limit int64)
	WriteMessage(messageType int, data []byte) error
}

type ClientList map[*Client]struct{}

type Client struct {
	connection WebSocketConnection
	manager    *Manager
	// the incoming information is grouped into a channel
	egress   chan Event
	nextPing *time.Ticker
}

func NewClient(conn WebSocketConnection, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

func (c *Client) readMessages() {
	defer c.manager.removeClient(c)
	// setup for receiving pong message
	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return
	}
	c.connection.SetPongHandler(c.pongHandler)
	// jumbo frame limit
	c.connection.SetReadLimit(1024)

	for {
		_, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error while reading the incoming message through ws: %v", err)
				return
			} else if err == websocket.ErrReadLimit {
				log.Printf("The incoming message was too big, connection closed: %v", err)
				return
			}
		}
		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error unmarshaling event: %v", err)
			return
		}

		// here, depending on the event type, we will call a function
		if err := c.manager.routeEvent(request, c); err != nil {
			return
		}
	}
}

func (c *Client) writeMessages() {
	defer c.manager.removeClient(c)
	// setup of receiving ping message
	c.nextPing = time.NewTicker(pingInterval)

	for {
		select {
		case event, ok := <-c.egress:
			// if there is no data in the channel but it was triggered, something went wrong
			// we close the connection
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}
			data, err := json.Marshal(event)
			if err != nil {
				log.Printf("failed to marshal the event: %v", err)
				return
			}
			err = c.connection.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("failed to marshal the event: %v", err)
				return
			}

		case <-c.nextPing.C:
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("writing err: %v", err)
				return
			}
			c.nextPing = time.NewTicker(pingInterval)
		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	// starting the ping
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}
