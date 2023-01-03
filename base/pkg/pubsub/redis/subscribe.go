package redis

import (
	"context"

	"github.com/go-redis/redis/v9"
)

func (p *redisPubSub) Subscribe(ctx context.Context) (<-chan []byte, error) {
	sub := p.redis.Subscribe(
		ctx,
		p.channel,
	).Channel(redis.WithChannelSize(64))
	ch := make(chan []byte, 32)
	go func() {
		defer close(ch)
		for {
			msg, ok := <-sub
			if !ok {
				return
			}
			select {
			case ch <- []byte(msg.Payload):
			default:
				p.log.Warn("redis pubsub dropped message due to channel being full")
			}
		}
	}()
	return ch, nil
}
