package redis

import (
	"context"

	"github.com/pkg/errors"
)

func (d *redisDataCache) Add(projectID string, videoID string) error {
	p := d.redis.Pipeline()
	ctx := context.Background()
	p.SAdd(ctx, videoProjectsKey(videoID), projectID)
	p.SAdd(ctx, datasetKey(projectID), videoID)
	if _, err := p.Exec(ctx); err != nil {
		return errors.Wrap(err, "redis")
	}
	return nil
}
