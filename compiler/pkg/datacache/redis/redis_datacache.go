package redis

import (
	"fmt"

	"github.com/go-redis/redis/v9"
	"github.com/thavlik/t4vd/compiler/pkg/datacache"
	"go.uber.org/zap"
)

type redisDataCache struct {
	redis *redis.Client
	log   *zap.Logger
}

func NewRedisDataCache(
	redis *redis.Client,
	log *zap.Logger,
) datacache.DataCache {
	return &redisDataCache{
		redis,
		log,
	}
}

func datasetKey(projectID string) string {
	return fmt.Sprintf("ds:%s", projectID)
}

func videoProjectsKey(videoID string) string {
	return fmt.Sprintf("vp:%s", videoID)
}
