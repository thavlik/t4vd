package redis

import (
	"context"

	"github.com/pkg/errors"
)

func (r *redisCachedSet) Complete(ctx context.Context, key string) error {
	p := r.redis.Pipeline()
	// set the complete flag to 1
	p.Set(
		ctx,
		completeKey(key),
		"1",
		0,
	)
	// push an empty value to the channel to signal completion
	p.Publish(ctx, membersKey(key), "")
	if _, err := p.Exec(ctx); err != nil {
		return errors.Wrap(err, "redis")
	}
	return nil
}
