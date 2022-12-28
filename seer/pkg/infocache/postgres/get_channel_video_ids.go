package postgres

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
)

func (c *postgresInfoCache) GetChannelVideoIDs(
	ctx context.Context,
	channelID string,
) ([]string, time.Time, error) {
	timestamp, err := getRecency(
		channelID,
		channelRecencyTable,
		c.db,
	)
	if err != nil {
		return nil, time.Time{}, err
	}
	videoIDs, err := getVideoIDs(
		ctx,
		channelJoinTable,
		channelOriginKey,
		channelID,
		c.db,
	)
	if err != nil {
		return nil, time.Time{}, errors.Wrap(err, "getVideoIDs")
	} else if len(videoIDs) == 0 {
		return nil, time.Time{}, infocache.ErrCacheUnavailable
	}
	return videoIDs, timestamp, nil
}
