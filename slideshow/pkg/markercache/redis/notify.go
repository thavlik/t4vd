package redis

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (m *redisMarkerCache) notify(projectID string) error {
	p := m.redis.Pipeline()
	notifyPipeline(context.Background(), p, projectID)
	if _, err := p.Exec(context.Background()); err != nil {
		return errors.Wrap(err, "redis.Exec")
	}
	m.log.Debug("notified queue", zap.String("projectID", projectID))
	return nil
}

func notifyPipeline(ctx context.Context, p redis.Pipeliner, projectID string) {
	p.ZIncrBy(ctx, queuedProjectsKey, 1.0, projectID) // add project queue or increase priority
	p.Publish(ctx, channelKey, projectID)             // notify queue has changed
}
