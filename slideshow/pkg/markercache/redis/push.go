package redis

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/slideshow/pkg/api"
	"go.uber.org/zap"
)

func (m *redisMarkerCache) Push(
	projectID string,
	markers ...*api.Marker,
) error {
	if len(markers) == 0 {
		return nil
	}
	key := projectQueueKey(projectID)
	p := m.redis.Pipeline()
	var r *redis.IntCmd
	for _, marker := range markers {
		body, err := json.Marshal(marker)
		if err != nil {
			return errors.Wrap(err, "json")
		}
		r = p.RPush(context.Background(), key, body)
	}
	if _, err := p.Exec(context.Background()); err != nil {
		return errors.Wrap(err, "redis pipeline")
	}
	numCached, _ := r.Result()
	if numCached >= m.remoteCacheCapacity {
		// remove the project from the queue
		if _, err := m.redis.ZRem(
			context.Background(),
			queuedProjectsKey,
			projectID,
		).Result(); err != nil {
			return errors.Wrap(err, "redis.ZRem")
		}
	} else {
		// deprioritize the project by the number of
		// markers we just pushed
		if _, err := m.redis.ZIncrBy(
			context.Background(),
			queuedProjectsKey,
			-1.0*float64(len(markers)),
			projectID,
		).Result(); err != nil {
			return errors.Wrap(err, "redis.ZIncrBy")
		}
	}
	m.log.Debug("pushed redis marker cache",
		zap.String("projectID", projectID),
		zap.Int("numPushed", len(markers)),
		zap.Int64("numCached", numCached))
	return nil
}
