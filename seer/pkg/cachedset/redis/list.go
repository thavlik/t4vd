package redis

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
)

func (r *redisCachedSet) List(
	ctx context.Context,
	key string,
) (values []string, complete bool, err error) {
	p := r.redis.Pipeline()
	members := p.ZRange(ctx, membersKey(key), 0, -1)
	isComplete := p.Get(ctx, completeKey(key))
	if _, err := p.Exec(ctx); err != nil && err != redis.Nil {
		return nil, false, errors.Wrap(err, "pipeline.Exec")
	}
	isComp, err := isComplete.Result()
	if err != nil && err != redis.Nil {
		return nil, false, errors.Wrap(err, "isComplete.Result")
	}
	complete = isComp == "1"
	values, err = members.Result()
	if err != nil && err != redis.Nil {
		return nil, false, errors.Wrap(err, "members.Result")
	}
	return values, complete, nil
}
