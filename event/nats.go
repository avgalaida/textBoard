package event

import (
	"bytes"
	"encoding/gob"
	"github.com/nats-io/go-nats"
	"textBoard/schema"
)

type NatsEventStore struct {
	nc                      *nats.Conn
	postCreatedSubscription *nats.Subscription
	postCreatedChan         chan PostCreatedMessage
}

func NewNats(url string) (*NatsEventStore, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsEventStore{nc: nc}, nil
}

func (e *NatsEventStore) Close() {
	if e.nc != nil {
		e.nc.Close()
	}
	if e.postCreatedSubscription != nil {
		e.postCreatedSubscription.Unsubscribe()
	}
	close(e.postCreatedChan)
}

func (e *NatsEventStore) PublishPostCreated(post schema.Post) error {
	m := PostCreatedMessage{post.ID, post.Body, post.CreatedAt}
	data, err := e.writeMessage(&m)
	if err != nil {
		return err
	}
	return e.nc.Publish(m.Key(), data)
}

func (mq *NatsEventStore) writeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
