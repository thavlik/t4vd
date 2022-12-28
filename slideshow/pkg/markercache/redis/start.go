package redis

func (m *redisMarkerCache) Start() {
	go m.worker()
}
