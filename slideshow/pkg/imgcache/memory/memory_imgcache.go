package memory

import (
	"context"
	"sync"
	"time"

	"github.com/thavlik/t4vd/slideshow/pkg/imgcache"
)

type cachedImg struct {
	data       []byte
	lastServed time.Time
}

type memoryImgCache struct {
	cap  int
	data map[string]*cachedImg
	l    sync.Mutex
}

func NewMemoryImgCache(capacity int) imgcache.ImgCache {
	return &memoryImgCache{
		cap:  capacity,
		data: make(map[string]*cachedImg),
	}
}

func (m *memoryImgCache) GetImage(
	ctx context.Context,
	videoID string,
	t time.Duration,
) ([]byte, error) {
	m.l.Lock()
	defer m.l.Unlock()
	img, ok := m.data[imgcache.MangleKey(videoID, t)]
	if !ok {
		return nil, imgcache.ErrNotCached
	}
	img.lastServed = time.Now()
	return img.data, nil
}

func (m *memoryImgCache) getLRU() string {
	var lru *time.Time
	var key string
	for k, v := range m.data {
		if lru == nil || v.lastServed.Before(*lru) {
			key = k
			lru = &v.lastServed
		}
	}
	return key
}

func (m *memoryImgCache) SetImage(
	videoID string,
	t time.Duration,
	img []byte,
) error {
	m.l.Lock()
	defer m.l.Unlock()
	if len(m.data) >= m.cap {
		delete(m.data, m.getLRU())
	}
	m.data[imgcache.MangleKey(videoID, t)] = &cachedImg{
		data:       img,
		lastServed: time.Now(),
	}
	return nil
}
