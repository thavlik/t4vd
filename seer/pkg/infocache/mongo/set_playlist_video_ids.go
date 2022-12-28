package mongo

import (
	"time"

	"github.com/pkg/errors"
)

func (c *mongoInfoCache) SetPlaylistVideoIDs(
	playlistID string,
	videoIDs []string,
	timestamp time.Time,
) error {
	if err := setVideoIDs(
		c.playlistJoinCollection,
		playlistOriginKey,
		playlistID,
		videoIDs,
	); err != nil {
		return errors.Wrap(err, "setVideoIDs")
	}
	if err := c.playlistRefreshed(playlistID, timestamp); err != nil {
		return errors.Wrap(err, "cachedPlaylistRefreshed")
	}
	return nil
}

func (c *mongoInfoCache) playlistRefreshed(
	playlistID string,
	timestamp time.Time,
) error {
	return refreshCache(
		c.playlistRecencyCollection,
		playlistID,
		timestamp,
	)
}
