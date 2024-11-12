package model

import (
	"time"

	"github.com/luuisavelino/network-interface/internal/domain/entities"
)

type RequestsResponse struct {
	Sent     []RequestResponse `json:"sent"`
	Received []RequestResponse `json:"received"`
}

func ToRequestsResponse(m entities.Requests) RequestsResponse {
	sent := make([]RequestResponse, 0)
	received := make([]RequestResponse, 0)

	for _, Request := range m.Sent {
		sent = append(sent, ToRequestResponse(*Request))
	}

	for _, Request := range m.Received {
		received = append(received, ToRequestResponse(*Request))
	}

	return RequestsResponse{
		Sent:     sent,
		Received: received,
	}
}

type RequestRequest struct {
	Header Header      `json:"header"`
	Body   interface{} `json:"body"`
}

type Header struct {
	Topic       string `json:"topic"`
	Sender      string `json:"sender"`
	Destination string `json:"destination"`
}

func (m RequestRequest) ToDomain() entities.Request {
	return entities.Request{
		Header: entities.Header{
			Sender:      m.Header.Sender,
			Destination: m.Header.Destination,
		},
		Body: m.Body,
	}
}

type RequestResponse struct {
	ID     string      `json:"id"`
	Header Header      `json:"header"`
	Body   interface{} `json:"content"`
	Read   bool        `json:"read"`
	Date   time.Time   `json:"date"`
}

func ToRequestResponse(m entities.Request) RequestResponse {
	var content string
	switch m.Header.Topic {
	case "user-message":
		content = m.Body.(string)
	default:
		content = m.Header.Topic
	}

	return RequestResponse{
		ID: m.ID.String(),
		Header: Header{
			Topic:       m.Header.Topic,
			Sender:      m.Header.Sender,
			Destination: m.Header.Destination,
		},
		Body: content,
		Read: m.IsRead(),
		Date: m.Date,
	}
}
