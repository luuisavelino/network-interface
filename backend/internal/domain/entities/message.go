package entities

import "github.com/google/uuid"

type Messages struct {
	All 	 map[uuid.UUID]Message
	Read 	 map[uuid.UUID]Message
	Unread map[uuid.UUID]Message
}

type Message struct {
	Topic   string
	Source  int
	Target  int
	Content interface{}
}

func NewMessage(topic string, source int, target int, content interface{}) Message {
	return Message{
		Topic: topic,
		Source: source,
		Target: target,
		Content: content,
	}
}
