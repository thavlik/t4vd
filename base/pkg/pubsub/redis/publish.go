package redis

import (
	"context"

	"github.com/pkg/errors"
)

func (p *redisPubSub) Publish(payload []byte) error {
	if _, err := p.redis.Publish(
		context.Background(),
		p.channel,
		payload,
	).Result(); err != nil {
		return errors.Wrap(err, "redis")
	}
	return nil
}
