package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/slideshow/pkg/imgcache"
)

type redisImgCache struct {
	redis  *redis.Client
	expiry time.Duration
}

func NewRedisImgCache(
	redis *redis.Client,
	expiry time.Duration,
) imgcache.ImgCache {
	return &redisImgCache{redis, expiry}
}

func (m *redisImgCache) GetImage(
	ctx context.Context,
	videoID string,
	t time.Duration,
) ([]byte, error) {
	v, err := m.redis.Get(
		ctx,
		imgcache.MangleKey(videoID, t),
	).Result()
	if err == redis.Nil {
		return nil, imgcache.ErrNotCached
	} else if err != nil {
		return nil, errors.Wrap(err, "redis.Get")
	}
	return []byte(v), nil
}

func (m *redisImgCache) SetImage(videoID string, t time.Duration, img []byte) error {
	if _, err := m.redis.Set(
		context.Background(),
		imgcache.MangleKey(videoID, t),
		img,
		m.expiry,
	).Result(); err != nil {
		return errors.Wrap(err, "redis.Set")
	}
	return nil
}
