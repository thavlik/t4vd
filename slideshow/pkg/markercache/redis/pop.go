package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/slideshow/pkg/api"
)

func (m *redisMarkerCache) Pop(
	ctx context.Context,
	projectID string,
) (*api.Marker, error) {
	qkey := projectQueueKey(projectID)
	done := ctx.Done()
	for {
		select {
		case <-done:
			return nil, ctx.Err()
		default:
			p := m.redis.Pipeline()
			lpop := p.LPop(ctx, qkey)         // pop the value
			notifyPipeline(ctx, p, projectID) // notify change
			if _, err := p.Exec(ctx); err != nil && err != redis.Nil {
				panic(errors.Wrap(err, "redis.Exec"))
			}
			value, err := lpop.Result()
			if err == redis.Nil {
				// no cached markers available, try again after waiting
				select {
				case <-done:
					return nil, ctx.Err()
				case <-time.After(delay):
					continue
				}
			} else if err != nil {
				panic(errors.Wrap(err, "redis.LPop"))
			}
			marker := &api.Marker{}
			if err := json.Unmarshal([]byte(value), marker); err != nil {
				panic(errors.Wrap(err, "json"))
			}
			return marker, nil
		}
	}
}
