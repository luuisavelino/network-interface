package entities

import (
	"time"

	"github.com/google/uuid"
)

type Requests struct {
	Sent     map[uuid.UUID]*Request
	Received map[uuid.UUID]*Request
}

type Request struct {
	ID     uuid.UUID
	Header Header
	Body   interface{}
	read   bool
	Date   time.Time
}

type Header struct {
	Topic       string
	Sender      string
	Destination string
	Path        []Route
	ContentType string
}

func NewRequest(topic, source, target string, path []Route, body interface{}) Request {
	return Request{
		ID: uuid.New(),
		Header: Header{
			Topic:       topic,
			Sender:      source,
			Destination: target,
			Path:        path,
		},
		Body: body,
		read: false,
		Date: time.Now(),
	}
}

func (m *Request) Read() {
	m.read = true
}

func (m *Request) IsRead() bool {
	return m.read
}
