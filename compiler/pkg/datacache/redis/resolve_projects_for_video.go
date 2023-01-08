package redis

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
)

func (d *redisDataCache) ResolveProjectsForVideo(
	ctx context.Context,
	videoID string,
) ([]string, error) {
	potential, err := d.redis.SMembers(
		ctx,
		videoProjectsKey(videoID),
	).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "redis")
	}
	var projectIDs []string
	for _, projectID := range potential {
		if isMember, err := d.redis.SIsMember(
			ctx,
			datasetKey(projectID),
			videoID,
		).Result(); err != nil {
			return nil, errors.Wrap(err, "redis.SIsMember")
		} else if isMember {
			projectIDs = append(projectIDs, projectID)
		} else {
			// remove the project from the set of potential projects
			if _, err := d.redis.SRem(
				ctx,
				videoProjectsKey(videoID),
				projectID,
			).Result(); err != nil {
				return nil, errors.Wrap(err, "redis.SRem")
			}
		}
	}
	return projectIDs, nil
}
