package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
)

func (c *mongoInfoCache) GetPlaylistVideoIDs(
	ctx context.Context,
	playlistID string,
) ([]string, time.Time, error) {
	timestamp, err := getRecency(c.playlistRecencyCollection, playlistID)
	if err != nil {
		return nil, time.Time{}, err
	}
	videoIDs, err := getVideoIDs(
		ctx,
		c.playlistJoinCollection,
		playlistOriginKey,
		playlistID,
	)
	if err != nil {
		return nil, time.Time{}, errors.Wrap(err, "getVideoIDs")
	} else if len(videoIDs) == 0 {
		return nil, time.Time{}, infocache.ErrCacheUnavailable
	}
	return videoIDs, timestamp, nil
}
