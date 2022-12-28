package postgres

import (
	"time"

	"github.com/pkg/errors"
)

func (c *postgresInfoCache) SetChannelVideoIDs(
	channelID string,
	videoIDs []string,
	timestamp time.Time,
) error {
	if err := setVideoIDs(
		channelJoinTable,
		channelOriginKey,
		channelID,
		videoIDs,
		c.db,
	); err != nil {
		return errors.Wrap(err, "setVideoIDs")
	}
	if err := c.channelRefreshed(channelID, timestamp); err != nil {
		return errors.Wrap(err, "channelRefreshed")
	}
	return nil
}

func (c *postgresInfoCache) channelRefreshed(
	channelID string,
	timestamp time.Time,
) error {
	return refreshCache(
		channelRecencyTable,
		channelID,
		timestamp,
		c.db,
	)
}
