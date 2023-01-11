package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/pubsub"
	"github.com/thavlik/t4vd/base/pkg/scheduler"
	hound "github.com/thavlik/t4vd/hound/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/cachedset"
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
		panic(base.Unreachable)
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
	cachedVideoIDs cachedset.CachedSet,
	pub pubsub.Publisher,
	querySched scheduler.Scheduler,
	hound hound.Hound,
	stop <-chan struct{},
	log *zap.Logger,
) {
	popQuery := make(chan string)
	go queryPopper(popQuery, querySched, stop, log)
	for i := 0; i < concurrency; i++ {
		go queryWorker(
			infoCache,
			thumbCache,
			popQuery,
			querySched,
			pub,
			cachedVideoIDs,
			hound,
			log,
		)
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
	pub pubsub.Publisher,
	cachedVideoIDs cachedset.CachedSet,
	houndClient hound.Hound,
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
		if err := func() error {
			defer lock.Release()
			ctx, cancel := context.WithCancel(context.Background())
			onProgress := make(chan struct{}, 4)
			stopped := make(chan struct{})
			defer func() {
				cancel()
				<-stopped
			}()
			go func() {
				defer func() { stopped <- struct{}{} }()
				for {
					select {
					case <-ctx.Done():
					case _, ok := <-onProgress:
						if !ok {
							return
						}
						if err := lock.Extend(); err != nil {
							panic(err)
						}
					}
				}
			}()
			switch e.Type {
			case playlist:
				if recent, err := infoCache.IsPlaylistRecent(e.ID); err != nil {
					return errors.Wrap(err, "infocache.IsPlaylistRecent")
				} else if !recent {
					playlist := &api.PlaylistDetails{}
					if err := queryPlaylist(e.ID, playlist); err != nil {
						return errors.Wrap(err, "failed to query playlist")
					}
					if err := infoCache.SetPlaylist(playlist); err != nil {
						return errors.Wrap(err, "infocache.SetPlaylist")
					}
					base.Progress(ctx, onProgress)
					go reportPlaylistDetails(ctx, playlist, houndClient, log)
					onVideo := make(chan *api.VideoDetails, 32)
					stopped := make(chan struct{})
					innerCtx, cancel := context.WithCancel(ctx)
					go func() {
						defer func() { stopped <- struct{}{} }()
						numVideos := 0
						for {
							select {
							case <-innerCtx.Done():
								return
							case video, ok := <-onVideo:
								if !ok {
									return
								}
								numVideos++
								reportPlaylistVideo(
									innerCtx,
									e.ID,
									video,
									numVideos,
									pub,
									cachedVideoIDs,
									houndClient,
									log,
								)
								base.Progress(innerCtx, onProgress)
							}
						}
					}()
					_, err := retrievePlaylistVideos(
						ctx,
						infoCache,
						e.ID,
						onVideo,
						entityLog,
					)
					cancel()
					if err != nil {
						return errors.Wrap(err, "failed to retrieve playlist videos")
					}
					<-stopped
					if err := cachedVideoIDs.Complete(
						ctx,
						e.ID,
					); err != nil {
						return errors.Wrap(err, "failed to complete cached video IDs")
					}
				}
				base.Progress(ctx, onProgress)
				if err := cachePlaylistThumbnail(
					ctx,
					e.ID,
					thumbCache,
					entityLog,
				); err != nil {
					return errors.Wrap(err, "failed to download playlist thumbnail")
				}
			case channel:
				var channel *api.ChannelDetails
				if recent, err := infoCache.IsChannelRecent(e.ID); err != nil {
					return errors.Wrap(err, "infocache.IsChannelRecent")
				} else if recent {
					channel, err = infoCache.GetChannel(ctx, e.ID)
					if err != nil {
						return errors.Wrap(err, "infocache.GetChannel")
					}
				} else {
					channel = &api.ChannelDetails{}
					if err := queryChannel(e.ID, channel); err != nil {
						return errors.Wrap(err, "failed to query channel")
					}
					if err := infoCache.SetChannel(channel); err != nil {
						return errors.Wrap(err, "infocache.SetChannel")
					}
					base.Progress(ctx, onProgress)
					go reportChannelDetails(ctx, channel, houndClient, log)
					onVideo := make(chan *api.VideoDetails, 32)
					stopped := make(chan struct{})
					innerCtx, cancel := context.WithCancel(ctx)
					go func() {
						defer func() { stopped <- struct{}{} }()
						numVideos := 0
						for {
							select {
							case <-innerCtx.Done():
								return
							case video, ok := <-onVideo:
								if !ok {
									return
								}
								numVideos++
								reportChannelVideo(
									innerCtx,
									e.ID,
									video,
									numVideos,
									pub,
									cachedVideoIDs,
									houndClient,
									log,
								)
								base.Progress(innerCtx, onProgress)
							}
						}
					}()
					_, err := retrieveChannelVideos(
						ctx,
						infoCache,
						e.ID,
						onVideo,
						entityLog,
					)
					cancel()
					if err != nil {
						return errors.Wrap(err, "failed to retrieve channel videos")
					}
					<-stopped
					if err := cachedVideoIDs.Complete(
						ctx,
						e.ID,
					); err != nil {
						return errors.Wrap(err, "failed to complete cached video IDs")
					}
				}
				base.Progress(ctx, onProgress)
				if err := cacheChannelAvatar(
					ctx,
					e.ID,
					channel.Avatar,
					thumbCache,
					entityLog,
				); err != nil {
					return errors.Wrap(err, "failed to download channel avatar")
				}
			case video:
				if recent, err := infoCache.IsVideoRecent(e.ID); err != nil {
					return errors.Wrap(err, "infocache.IsVideoRecent")
				} else if !recent {
					base.Progress(ctx, onProgress)
					video, err := queryVideoDetails(e.ID, log)
					if err != nil {
						return errors.Wrap(err, "queryVideoDetails")
					}
					base.Progress(ctx, onProgress)
					if err := infoCache.SetVideo(video); err != nil {
						return errors.Wrap(err, "infocache.SetVideo")
					}
					go reportVideoDetails(ctx, video, houndClient, log)
				}
				base.Progress(ctx, onProgress)
				if err := cacheVideoThumbnail(
					ctx,
					e.ID,
					thumbCache,
					entityLog,
				); err != nil {
					return errors.Wrap(err, "failed to download video thumbnail")
				}
			default:
				panic(base.Unreachable)
			}
			if err := querySched.Remove(rawEnt); err != nil {
				return errors.Wrap(err, "failed to remove entity from query scheduler, this will result in multiple repeated requests to youtube")
			}
			return nil
		}(); err != nil {
			entityLog.Warn("query worker error", zap.Error(err))
		}
	}
}

func reportVideoDetails(
	ctx context.Context,
	details *api.VideoDetails,
	houndClient hound.Hound,
	log *zap.Logger,
) {
	if houndClient == nil {
		return
	}
	if _, err := houndClient.ReportVideoDetails(
		ctx,
		*((*hound.VideoDetails)(details)),
	); err != nil {
		log.Warn("failed to report video details", zap.Error(err))
	}
}

func reportPlaylistVideo(
	ctx context.Context,
	playlistID string,
	video *api.VideoDetails,
	numVideos int,
	pub pubsub.Publisher,
	cachedVideoIDs cachedset.CachedSet,
	houndClient hound.Hound,
	log *zap.Logger,
) {
	if err := pub.Publish(
		ctx,
		playlistTopic(playlistID),
		[]byte(video.ID),
	); err != nil {
		log.Warn("failed to publish video ID for playlist",
			zap.Error(err))
	}
	if err := cachedVideoIDs.Set(
		ctx,
		playlistID,
		video.ID,
		numVideos-1,
	); err != nil {
		log.Warn("failed to set cached video ID for playlist",
			zap.Error(err))
	}
	if houndClient != nil {
		if _, err := houndClient.ReportPlaylistVideo(
			ctx,
			hound.PlaylistVideo{
				PlaylistID: playlistID,
				Video:      *(*hound.VideoDetails)(video),
				NumVideos:  numVideos,
			},
		); err != nil {
			log.Warn("failed to report playlist video", zap.Error(err))
		}
	}
}

func reportChannelVideo(
	ctx context.Context,
	channelID string,
	video *api.VideoDetails,
	numVideos int,
	pub pubsub.Publisher,
	cachedVideoIDs cachedset.CachedSet,
	houndClient hound.Hound,
	log *zap.Logger,
) {
	if err := pub.Publish(
		ctx,
		channelTopic(channelID),
		[]byte(video.ID),
	); err != nil {
		log.Warn("failed to publish video ID for channel",
			zap.Error(err))
	}
	if err := cachedVideoIDs.Set(
		ctx,
		channelID,
		video.ID,
		numVideos-1,
	); err != nil {
		log.Warn("failed to set cached video ID for channel",
			zap.Error(err))
	}
	if houndClient != nil {
		if _, err := houndClient.ReportChannelVideo(
			ctx,
			hound.ChannelVideo{
				ChannelID: channelID,
				NumVideos: numVideos,
				Video:     *(*hound.VideoDetails)(video),
			},
		); err != nil {
			log.Warn("failed to report channel video", zap.Error(err))
		}
	}
}

func reportChannelDetails(
	ctx context.Context,
	details *api.ChannelDetails,
	houndClient hound.Hound,
	log *zap.Logger,
) {
	if houndClient == nil {
		return
	}
	if _, err := houndClient.ReportChannelDetails(
		ctx,
		hound.ChannelDetails{
			ID:     details.ID,
			Name:   details.Name,
			Avatar: details.Avatar,
			Subs:   details.Subs,
		},
	); err != nil {
		log.Warn("failed to report channel details", zap.Error(err))
	}
}

func reportPlaylistDetails(
	ctx context.Context,
	details *api.PlaylistDetails,
	houndClient hound.Hound,
	log *zap.Logger,
) {
	if houndClient == nil {
		return
	}
	if _, err := houndClient.ReportPlaylistDetails(
		ctx,
		hound.PlaylistDetails{
			ID:        details.ID,
			Title:     details.Title,
			Channel:   details.Channel,
			ChannelID: details.ChannelID,
			NumVideos: details.NumVideos,
		},
	); err != nil {
		log.Warn("failed to report playlist details", zap.Error(err))
	}
}
