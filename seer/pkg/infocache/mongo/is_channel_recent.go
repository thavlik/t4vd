package mongo

func (c *mongoInfoCache) IsChannelRecent(
	channelID string,
) (bool, error) {
	return checkCacheRecency(c.channelRecencyCollection, channelID)
}
