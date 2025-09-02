package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type WebSocketConnection interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}

type ClientList map[*Client]struct{}

type Client struct {
	connection WebSocketConnection
	manager    *Manager
	// the incoming information is grouped into a channel
	egress chan Event
}

func NewClient(conn WebSocketConnection, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

func (c *Client) readMessages() error {
	for {
		_, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error while reading the incoming message through ws: %v", err)
				return err
			}
		}
		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error unmarshaling event: %v", err)
			return err
		}

		// here, depending on the event type, we will call a function
		if err := c.manager.routeEvent(request, c); err != nil {
			return err
		}
	}
}

func (c *Client) writeMessages() {
	for {
		event, ok := <-c.egress
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
	}
}
