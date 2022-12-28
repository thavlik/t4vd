package postgres

import (
	"time"

	"github.com/pkg/errors"
)

func (c *postgresInfoCache) SetPlaylistVideoIDs(
	playlistID string,
	videoIDs []string,
	timestamp time.Time,
) error {
	if err := setVideoIDs(
		playlistJoinTable,
		playlistOriginKey,
		playlistID,
		videoIDs,
		c.db,
	); err != nil {
		return errors.Wrap(err, "setVideoIDs")
	}
	if err := c.playlistRefreshed(playlistID, timestamp); err != nil {
		return errors.Wrap(err, "cachedPlaylistRefreshed")
	}
	return nil
}

func (c *postgresInfoCache) playlistRefreshed(
	playlistID string,
	timestamp time.Time,
) error {
	return refreshCache(
		playlistRecencyTable,
		playlistID,
		timestamp,
		c.db,
	)
}
