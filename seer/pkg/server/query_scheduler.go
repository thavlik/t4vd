package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/scheduler"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	"github.com/thavlik/t4vd/seer/pkg/thumbcache"
	"go.uber.org/zap"
)

type entityType int

func (e entityType) String() string {
	switch e {
	case playlist:
		return "playlist"
	case channel:
		return "channel"
	case video:
		return "video"
	default:
		panic("unreachable branch detected")
	}
}

const (
	playlist entityType = 1
	channel  entityType = 2
	video    entityType = 3
)

type entity struct {
	Type entityType `json:"t"`
	ID   string     `json:"_"`
}

func initQueryWorkers(
	concurrency int,
	infoCache infocache.InfoCache,
	thumbCache thumbcache.ThumbCache,
	querySched scheduler.Scheduler,
	stop <-chan struct{},
	log *zap.Logger,
) {
	popQuery := make(chan string)
	go queryPopper(popQuery, querySched, stop, log)
	for i := 0; i < concurrency; i++ {
		go queryWorker(infoCache, thumbCache, popQuery, querySched, log)
	}
}

func queryPopper(
	popQuery chan<- string,
	querySched scheduler.Scheduler,
	stop <-chan struct{},
	log *zap.Logger,
) {
	notification := querySched.Notify()
	defer close(popQuery)
	delay := 12 * time.Second
	for {
		start := time.Now()
		entities, err := querySched.List()
		if err != nil {
			panic(errors.Wrap(err, "scheduler.List"))
		}
		if len(entities) > 0 {
			log.Debug("checking entities", zap.Int("num", len(entities)))
			for _, entJson := range entities {
				popQuery <- entJson
			}
		}
		remaining := delay - time.Since(start)
		if remaining < 0 {
			remaining = time.Millisecond
		}
		select {
		case <-stop:
			return
		case <-notification:
			continue
		case <-time.After(remaining):
			continue
		}
	}
}

func queryWorker(
	infoCache infocache.InfoCache,
	thumbCache thumbcache.ThumbCache,
	popQuery <-chan string,
	querySched scheduler.Scheduler,
	log *zap.Logger,
) {
	for {
		rawEnt, ok := <-popQuery
		if !ok {
			return
		}
		e := &entity{}
		if err := json.Unmarshal([]byte(rawEnt), e); err != nil {
			fmt.Println("removing malformed rawEnt")
			fmt.Println(string(rawEnt))
			if err := querySched.Remove(rawEnt); err != nil {
				panic(err)
			}
			continue
		}
		entityLog := log.With(
			zap.String("type", e.Type.String()),
			zap.String("id", e.ID))
		lock, err := querySched.Lock(e.ID)
		if err == scheduler.ErrLocked {
			// go to the next project
			entityLog.Debug("entity already locked")
			continue
		} else if err != nil {
			panic(errors.Wrap(err, "scheduler.Lock"))
		}
		entityLog.Debug("processing query item")
		func() {
			defer lock.Release()
			gotVideo := make(chan *api.VideoDetails)
			go func() {
				for {
					_, ok := <-gotVideo
					if !ok {
						return
					}
					if err := lock.Extend(); err != nil {
						panic(err)
					}
				}
			}()
			switch e.Type {
			case playlist:
				if err := cachePlaylistThumbnail(
					context.Background(),
					e.ID,
					thumbCache,
					entityLog,
				); err != nil {
					entityLog.Warn("failed to download video thumbnail", zap.Error(err))
				}
				if recent, err := infoCache.IsPlaylistRecent(e.ID); err != nil {
					entityLog.Warn("infocache.IsPlaylistRecent", zap.Error(err))
				} else if !recent {
					if _, err := retrievePlaylistVideos(
						infoCache,
						e.ID,
						gotVideo,
						entityLog,
					); err != nil {
						entityLog.Warn("failed to retrieve playlist videos", zap.Error(err))
					}
				}
			case channel:
				details, err := infoCache.GetChannel(context.Background(), e.ID)
				if err != nil {
					entityLog.Warn("failed to retrieve playlist details", zap.Error(err))
					return
				}
				if err := cacheChannelAvatar(
					context.Background(),
					e.ID,
					details.Avatar,
					thumbCache,
					entityLog,
				); err != nil {
					entityLog.Warn("failed to download video thumbnail", zap.Error(err))
				}
				if recent, err := infoCache.IsChannelRecent(e.ID); err != nil {
					entityLog.Warn("infocache.IsChannelRecent", zap.Error(err))
				} else if !recent {
					if _, err := retrieveChannelVideos(
						infoCache,
						e.ID,
						gotVideo,
						entityLog,
					); err != nil {
						entityLog.Warn("failed to retrieve channel videos", zap.Error(err))
					}
				}
			case video:
				if err := cacheVideoThumbnail(
					context.Background(),
					e.ID,
					thumbCache,
					entityLog,
				); err != nil {
					entityLog.Warn("failed to download video thumbnail", zap.Error(err))
				}
			default:
				panic("unreachable branch detected")
			}
			if err := querySched.Remove(rawEnt); err != nil {
				entityLog.Warn("failed to remove entity from query scheduler, this will result in multiple repeated requests to youtube")
			} else {
				entityLog.Debug("removed entity from query schedule")
			}
		}()
	}
}
