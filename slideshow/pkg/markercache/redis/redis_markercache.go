package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/thavlik/t4vd/slideshow/pkg/markercache"
	"go.uber.org/zap"
)

var (
	channelKey        = "sldshw"
	queuedProjectsKey = "sldshwq"
	delay             = 317 * time.Millisecond
)

type redisMarkerCache struct {
	remoteCacheCapacity int64
	redis               *redis.Client
	log                 *zap.Logger
	genMarker           markercache.GenMarkerFunc
	stop                chan struct{}
	stopped             chan struct{}
}

func NewRedisMarkerCache(
	redisClient *redis.Client,
	genMarker markercache.GenMarkerFunc,
	remoteCacheCapacity int64,
	log *zap.Logger,
) markercache.MarkerCache {
	m := &redisMarkerCache{
		remoteCacheCapacity: remoteCacheCapacity,
		redis:               redisClient,
		log:                 log,
		genMarker:           genMarker,
		stop:                make(chan struct{}, 1),
		stopped:             make(chan struct{}),
	}
	return m
}

func projectQueueKey(projectID string) string {
	return fmt.Sprintf("sldshw:{%s}", projectID)
}
