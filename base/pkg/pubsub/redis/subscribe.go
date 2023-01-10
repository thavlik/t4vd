package redis

import (
	"context"

	"github.com/thavlik/t4vd/base/pkg/pubsub"
)

func (p *redisPubSub) Subscribe(
	ctx context.Context,
	topic string,
) (pubsub.Subscription, error) {
	return &redisSubscription{
		redis: p.redis,
		stop:  make(chan struct{}, 1),
		topic: topic,
		log:   p.log,
	}, nil
}
