package redis

import (
	"github.com/go-redis/redis/v9"
	"github.com/thavlik/t4vd/base/pkg/pubsub"
	"go.uber.org/zap"
)

type redisPubSub struct {
	redis   *redis.Client
	channel string
	log     *zap.Logger
}

func NewRedisPubSub(
	redis *redis.Client,
	channel string,
	log *zap.Logger,
) pubsub.PubSub {
	return &redisPubSub{
		redis,
		channel,
		log,
	}
}
