package ws

import (
	"log"

	"github.com/google/uuid"
)

type Event struct {
	Type    string    `json:"type"`
	Room    uuid.UUID `json:"room"`
	Payload []byte    `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

type EventRelation map[string]EventHandler

const (
	EventSendMessage string = "send_message"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

func SendMessage(event Event, c *Client) error {
	// send the information to the users
	// for now we broadcast the message to all users...

	log.Println(event.Room, string(event.Payload))
	for client := range c.manager.clients {
		if client == c {
			continue
		}
		client.egress <- event
	}
	return nil
}

func getEventRelations() EventRelation {
	er := make(map[string]EventHandler)
	er["send_message"] = SendMessage

	return er
}
