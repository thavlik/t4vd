package mongo

func (c *mongoInfoCache) IsVideoRecent(
	videoID string,
) (bool, error) {
	return checkCacheRecency(c.playlistRecencyCollection, videoID)
}
