package pubsub

import "context"

type Publisher interface {
	Publish(message []byte) error
}

type Subscriber interface {
	Subscribe(context.Context) (<-chan []byte, error)
}

type PubSub interface {
	Publisher
	Subscriber
}
