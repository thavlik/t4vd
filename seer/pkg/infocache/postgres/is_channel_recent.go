package postgres

func (c *postgresInfoCache) IsChannelRecent(
	channelID string,
) (bool, error) {
	return checkCacheRecency(channelID, channelRecencyTable, c.db)
}
