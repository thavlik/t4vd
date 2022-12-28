package redis

import (
	"context"
	"strings"

	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/compiler/pkg/datastore"
	"go.uber.org/zap"
)

// populate generates a single marker and adds
// it to the project's queue
func (m *redisMarkerCache) populate(
	projectID string,
) error {
	m.log.Debug("populating queue", zap.String("projectID", projectID))
	marker, err := m.genMarker(projectID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.ErrDatasetNotFound.Error()) {
			// This happens when there aren't any videos
			// in a project or before a project has been
			// compiled to any extent
			if _, err := m.redis.ZRem(
				context.Background(),
				queuedProjectsKey,
				projectID,
			).Result(); err != nil && err != redis.Nil {
				return errors.Wrap(err, "redis.ZRem")
			}
			return nil
		}
		return errors.Wrap(err, "genMarker")
	}
	m.log.Debug("pushing queue", zap.String("projectID", projectID))
	return m.Push(projectID, marker)
}
