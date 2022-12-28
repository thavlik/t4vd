package redis

import (
	"github.com/pkg/errors"
)

func (m *redisMarkerCache) Close() error {
	m.stop <- struct{}{}
	<-m.stopped
	if err := m.redis.Close(); err != nil {
		return errors.Wrap(err, "redis.Close")
	}
	return nil
}
