package mongo

func (c *mongoInfoCache) IsPlaylistRecent(
	playlistID string,
) (bool, error) {
	return checkCacheRecency(c.playlistRecencyCollection, playlistID)
}
