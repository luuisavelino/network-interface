package model

import (
	"github.com/luuisavelino/network-interface/internal/domain/entities"
	"time"
)

type MessagesResponse struct {
	Sent     []MessageResponse `json:"sent"`
	Received []MessageResponse `json:"received"`
}

func ToMessagesResponse(m entities.Messages) MessagesResponse {
	sent := make([]MessageResponse, 0)
	received := make([]MessageResponse, 0)

	for _, message := range m.Sent {
		sent = append(sent, ToMessageResponse(*message))
	}

	for _, message := range m.Received {
		received = append(received, ToMessageResponse(*message))
	}

	return MessagesResponse{
		Sent:     sent,
		Received: received,
	}
}

type MessageResponse struct {
	ID      string      `json:"id"`
	Topic   string      `json:"topic"`
	Sender  int         `json:"sender"`
	Destination  int    `json:"destination"`
	Content interface{} `json:"content"`
	Read    bool        `json:"read"`
	Date 	time.Time     `json:"date"`
}

func ToMessageResponse(m entities.Message) MessageResponse {
	var content string
	switch m.Topic {
	case  "update-routing" :
		content = "Update routing table"
	default:
		content = m.Content.(string)
	}

	return MessageResponse{
		ID: m.ID.String(),
		Topic: m.Topic,
		Sender: m.Sender,
		Destination: m.Destination,
		Content: content,
		Read: m.IsRead(),
		Date: m.Date,
	}
}
