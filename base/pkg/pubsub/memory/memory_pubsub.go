package redis

import (
	"sync"

	"github.com/thavlik/t4vd/base/pkg/pubsub"
	"go.uber.org/zap"
)

type memoryPubSub struct {
	l        sync.Mutex
	channels map[chan []byte]struct{}
	log      *zap.Logger
}

func NewMemoryPubSub(log *zap.Logger) pubsub.PubSub {
	return &memoryPubSub{
		channels: make(map[chan []byte]struct{}),
		log:      log,
	}
}
