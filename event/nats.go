package event

import (
	"bytes"
	"encoding/gob"
	"github.com/avgalaida/textBoard/schema"
	"github.com/nats-io/nats.go"
	"log"
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

func (es *NatsEventStore) SubscribePostCreated() (<-chan PostCreatedMessage, error) {
	m := PostCreatedMessage{}
	es.postCreatedChan = make(chan PostCreatedMessage, 64)
	ch := make(chan *nats.Msg, 64)
	var err error
	es.postCreatedSubscription, err = es.nc.ChanSubscribe(m.Key(), ch)
	if err != nil {
		return nil, err
	}
	// Decode message
	go func() {
		for {
			select {
			case msg := <-ch:
				if err := es.readMessage(msg.Data, &m); err != nil {
					log.Fatal(err)
				}
				es.postCreatedChan <- m
			}
		}
	}()
	return (<-chan PostCreatedMessage)(es.postCreatedChan), nil
}

func (es *NatsEventStore) OnPostCreated(f func(PostCreatedMessage)) (err error) {
	m := PostCreatedMessage{}
	es.postCreatedSubscription, err = es.nc.Subscribe(m.Key(), func(msg *nats.Msg) {
		if err := es.readMessage(msg.Data, &m); err != nil {
			log.Fatal(err)
		}
		f(m)
	})
	return
}

func (es *NatsEventStore) Close() {
	if es.nc != nil {
		es.nc.Close()
	}
	if es.postCreatedSubscription != nil {
		if err := es.postCreatedSubscription.Unsubscribe(); err != nil {
			log.Fatal(err)
		}
	}
	close(es.postCreatedChan)
}

func (es *NatsEventStore) PublishPostCreated(post schema.Post) error {
	m := PostCreatedMessage{post.ID, post.Body, post.CreatedAt}
	data, err := es.writeMessage(&m)
	if err != nil {
		return err
	}
	return es.nc.Publish(m.Key(), data)
}

func (es *NatsEventStore) writeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (es *NatsEventStore) readMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}
