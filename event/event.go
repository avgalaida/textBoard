package event

import "github.com/avgalaida/textBoard/schema"

type Store interface {
	Close()
	PublishPostCreated(post schema.Post) error
	SubscribePostCreated() (<-chan PostCreatedMessage, error)
	OnPostCreated(f func(PostCreatedMessage)) error
}

var impl Store

func SetEventStore(es Store) {
	impl = es
}

func Close() {
	impl.Close()
}

func PublishPostCreated(post schema.Post) error {
	return impl.PublishPostCreated(post)
}

func SubscribePostCreated() (<-chan PostCreatedMessage, error) {
	return impl.SubscribePostCreated()
}

func OnPostCreated(f func(PostCreatedMessage)) error {
	return impl.OnPostCreated(f)
}
