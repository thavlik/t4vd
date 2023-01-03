package redis

import (
	"context"

	"github.com/pkg/errors"
)

func (d *redisDataCache) Set(projectID string, videoIDs []string) error {
	p := d.redis.Pipeline()
	key := datasetKey(projectID)
	p.Del(context.Background(), key)
	members := make([]interface{}, len(videoIDs))
	for i, videoID := range videoIDs {
		members[i] = videoID
		p.SAdd(
			context.Background(),
			videoProjectsKey(videoID),
			projectID,
		)
	}
	p.SAdd(context.Background(), key, members...)
	if _, err := p.Exec(context.TODO()); err != nil {
		return errors.Wrap(err, "redis")
	}
	return nil
}
