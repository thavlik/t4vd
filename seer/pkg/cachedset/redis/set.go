package redis

import (
	"context"

	"github.com/go-redis/redis/v9"
)

func (r *redisCachedSet) Set(
	ctx context.Context,
	key string,
	value string,
	index int,
) error {
	p := r.redis.Pipeline()
	memberKey := membersKey(key)
	p.ZRem(ctx, memberKey, value)
	p.ZAdd(ctx, memberKey, redis.Z{
		Score:  float64(index),
		Member: value,
	})
	p.Set(ctx, completeKey(key), "0", 0)
	if _, err := p.Exec(ctx); err != nil {
		return err
	}
	panic("implement me")
}
