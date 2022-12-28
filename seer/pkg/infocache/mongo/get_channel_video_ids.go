package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
)

func (c *mongoInfoCache) GetChannelVideoIDs(
	ctx context.Context,
	channelID string,
) ([]string, time.Time, error) {
	timestamp, err := getRecency(c.channelRecencyCollection, channelID)
	if err != nil {
		return nil, time.Time{}, err
	}
	videoIDs, err := getVideoIDs(
		ctx,
		c.channelJoinCollection,
		channelOriginKey,
		channelID,
	)
	if err != nil {
		return nil, time.Time{}, errors.Wrap(err, "getVideoIDs")
	} else if len(videoIDs) == 0 {
		return nil, time.Time{}, infocache.ErrCacheUnavailable
	}
	return videoIDs, timestamp, nil
}
