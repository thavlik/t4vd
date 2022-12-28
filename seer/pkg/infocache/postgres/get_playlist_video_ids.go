package postgres

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/seer/pkg/infocache"
)

func (c *postgresInfoCache) GetPlaylistVideoIDs(
	ctx context.Context,
	playlistID string,
) ([]string, time.Time, error) {
	timestamp, err := getRecency(
		playlistID,
		playlistRecencyTable,
		c.db,
	)
	if err != nil {
		return nil, time.Time{}, err
	}
	videoIDs, err := getVideoIDs(
		ctx,
		playlistJoinTable,
		playlistOriginKey,
		playlistID,
		c.db,
	)
	if err != nil {
		return nil, time.Time{}, errors.Wrap(err, "getVideoIDs")
	} else if len(videoIDs) == 0 {
		return nil, time.Time{}, infocache.ErrCacheUnavailable
	}
	return videoIDs, timestamp, nil
}
