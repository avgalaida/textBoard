package event

import "time"

type Message interface {
	Key() string
}

type PostCreatedMessage struct {
	ID        string
	Body      string
	CreatedAt time.Time
}

func (m *PostCreatedMessage) Key() string {
	return "post.created"
}
