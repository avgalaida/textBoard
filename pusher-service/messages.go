package main

import "time"

const (
	KindPostCreated = iota + 1
)

type PostCreatedMessage struct {
	Kind      uint32    `json:"kind"`
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

func newPostCreatedMessage(id string, body string, createdAt time.Time) *PostCreatedMessage {
	return &PostCreatedMessage{
		Kind:      KindPostCreated,
		ID:        id,
		Body:      body,
		CreatedAt: createdAt,
	}
}
