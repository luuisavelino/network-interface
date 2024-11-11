package entities

import (
	"time"

	"github.com/google/uuid"
)

type Messages struct {
	Sent     map[uuid.UUID]*Message
	Received map[uuid.UUID]*Message
}

type Message struct {
	ID          uuid.UUID
	Topic       string
	Sender      string
	Destination string
	Content     interface{}
	read        bool
	Date        time.Time
}

func NewMessage(topic string, source string, target string, content interface{}) Message {
	return Message{
		ID:          uuid.New(),
		Topic:       topic,
		Sender:      source,
		Destination: target,
		Content:     content,
		read:        false,
		Date:        time.Now(),
	}
}

func (m *Message) Read() {
	m.read = true
}

func (m *Message) IsRead() bool {
	return m.read
}
