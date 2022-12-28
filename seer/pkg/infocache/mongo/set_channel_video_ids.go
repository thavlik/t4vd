package mongo

import (
	"time"

	"github.com/pkg/errors"
)

func (c *mongoInfoCache) SetChannelVideoIDs(
	channelID string,
	videoIDs []string,
	timestamp time.Time,
) error {
	if err := setVideoIDs(
		c.channelJoinCollection,
		channelOriginKey,
		channelID,
		videoIDs,
	); err != nil {
		return errors.Wrap(err, "setVideoIDs")
	}
	if err := c.channelRefreshed(channelID, timestamp); err != nil {
		return errors.Wrap(err, "channelRefreshed")
	}
	return nil
}

func (c *mongoInfoCache) channelRefreshed(
	channelID string,
	timestamp time.Time,
) error {
	return refreshCache(
		c.channelRecencyCollection,
		channelID,
		timestamp,
	)
}
