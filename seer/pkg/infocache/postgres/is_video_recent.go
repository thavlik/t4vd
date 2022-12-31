package postgres

func (c *postgresInfoCache) IsVideoRecent(
	videoID string,
) (bool, error) {
	return checkCacheRecency(videoID, videoRecencyTable, c.db)
}
