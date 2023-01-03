package redis

import (
	"context"

	"github.com/go-redis/redis/v9"

	"github.com/pkg/errors"
)

func (d *redisDataCache) Get(
	ctx context.Context,
	projectID string,
) ([]string, error) {
	videoIDs, err := d.redis.SMembers(
		ctx,
		datasetKey(projectID),
	).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "redis")
	}
	return videoIDs, nil
}
