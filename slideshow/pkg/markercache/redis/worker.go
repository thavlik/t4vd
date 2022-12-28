package redis

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/thavlik/t4vd/base/pkg/base"
	"go.uber.org/zap"
)

// worker a goroutine that listens for events
// on a redis channel to know when to begin
// pushing to the buffer again
func (m *redisMarkerCache) worker() {
	m.log.Debug("redis marker cache worker thread started")
	defer func() { m.stopped <- struct{}{} }()
	ch := m.redis.Subscribe(context.Background(), channelKey).Channel()
	for {
		select {
		case <-m.stop:
			return
		case msg, ok := <-ch:
			if !ok {
				// redis connection terminated
				return
			} else if msg.Channel != channelKey {
				panic(base.Unreachable)
			}
			projectID := msg.Payload
			m.log.Debug("received notification", zap.String("projectID", projectID))
			if err := m.populate(projectID); err != nil {
				m.log.Warn("failed to populate project cache", zap.Error(err))
				continue
			}
		default:
			// if there are no ongoing notifications,
			// default to checking the queue of projects
			// that may not have full caches
			// the project getting hammered the hardest
			// should have the highest score
			result, err := m.redis.ZRevRangeByScore(
				context.Background(),
				queuedProjectsKey,
				&redis.ZRangeBy{
					Min:    "-inf",
					Max:    "+inf",
					Offset: 0,
					Count:  1,
				},
			).Result()
			if err != nil && err != redis.Nil {
				panic(err)
			} else if len(result) == 0 {
				continue
			}
			// populate only the first then continue
			// through this way, all project queues
			// will eventually be filled, but projects
			// that are being hammered should be first
			projectID := result[0]
			m.log.Debug("populating unfilled cache", zap.String("projectID", projectID))
			if err := m.populate(projectID); err != nil {
				m.log.Warn("failed to populate project cache", zap.Error(err))
				continue
			}
		}
	}
}
