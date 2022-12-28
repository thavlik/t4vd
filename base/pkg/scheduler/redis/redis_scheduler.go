package redis

import (
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v9"
	"github.com/thavlik/t4vd/base/pkg/scheduler"
)

type redisScheduler struct {
	redis  *redis.Client
	locker *redislock.Client
	key    string
	ttl    time.Duration
}

func NewRedisScheduler(
	redisClient *redis.Client,
	key string,
	ttl time.Duration,
) scheduler.Scheduler {
	return &redisScheduler{
		redisClient,
		redislock.New(redisClient),
		key,
		ttl,
	}
}

// channelName redis publish/subscribe channel name
func channelName(key string) string {
	return key + ":ch"
}
