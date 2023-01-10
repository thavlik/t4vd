package memory

import (
	"github.com/thavlik/t4vd/base/pkg/pubsub"
	"go.uber.org/zap"
)

type memoryPubSub struct {
	l        chan struct{}
	channels map[string]map[*memorySubscription]struct{}
	log      *zap.Logger
}

func NewMemoryPubSub(log *zap.Logger) pubsub.PubSub {
	return &memoryPubSub{
		l:        make(chan struct{}, 1),
		channels: make(map[string]map[*memorySubscription]struct{}),
		log:      log,
	}
}
