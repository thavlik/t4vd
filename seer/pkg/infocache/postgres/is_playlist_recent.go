package postgres

func (c *postgresInfoCache) IsPlaylistRecent(
	playlistID string,
) (bool, error) {
	return checkCacheRecency(playlistID, playlistRecencyTable, c.db)
}
