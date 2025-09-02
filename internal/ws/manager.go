package ws

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/CTSDM/motbwa-be/internal/api"
	"github.com/gorilla/websocket"
)

var webSocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Manager struct {
	clients  ClientList
	handlers map[string]EventHandler
}

func NewManager(ctx context.Context) *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}
	m.handlers = getEventRelations()
	return m
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("there is no such event type: %s", event.Type)
	}
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	// for now we assume there is no login
	log.Println("New websocket connection")

	conn, err := webSocketUpgrader.Upgrade(w, r, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		log.Printf("Couldn't update the websocket connection: %v", err)
		api.RespondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	} else if err != nil {
		log.Printf("Couldn't update the websocket connection: %v", err)
		api.RespondWithError(w, http.StatusInternalServerError, "Something went wrong while upgrading the websocket connection", err)
		return
	}

	client := NewClient(conn, m)
	m.clients[client] = struct{}{}

	go client.readMessages()
	go client.writeMessages()
}
