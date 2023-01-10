package redis

import (
	"fmt"

	"github.com/go-redis/redis/v9"
	"github.com/thavlik/t4vd/seer/pkg/cachedset"
	"go.uber.org/zap"
)

type redisCachedSet struct {
	redis *redis.Client
	log   *zap.Logger
}

func NewRedisCachedSet(
	redis *redis.Client,
	log *zap.Logger,
) cachedset.CachedSet {
	return &redisCachedSet{redis, log}
}

func completeKey(id string) string {
	return fmt.Sprintf("csc:%s", id)
}

func membersKey(id string) string {
	return fmt.Sprintf("csm:%s", id)
}
