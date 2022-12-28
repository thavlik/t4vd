package redis

func (m *redisMarkerCache) Queue(
	projectID string,
) error {
	return m.notify(projectID)
}
