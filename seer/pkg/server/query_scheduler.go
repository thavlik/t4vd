package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/scheduler"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	"github.com/thavlik/t4vd/seer/pkg/thumbcache"
	sources "github.com/thavlik/t4vd/sources/pkg/api"
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
	querySched scheduler.Scheduler,
	sources sources.Sources,
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
			sources,
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
	sourcesClient sources.Sources,
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
			done := ctx.Done()
			defer cancel()
			onProgress := make(chan struct{}, 4)
			stopped := make(chan struct{})
			defer func() {
				close(onProgress)
				<-stopped
			}()
			go func() {
				defer func() { stopped <- struct{}{} }()
				for {
					select {
					case <-done:
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
					go reportPlaylistDetails(playlist, sourcesClient, log)
					onVideo := make(chan *api.VideoDetails, 4)
					go func() {
						numVideos := 0
						for {
							select {
							case <-done:
								return
							case video, ok := <-onVideo:
								if !ok {
									return
								}
								numVideos++
								go reportPlaylistVideo(
									e.ID,
									video,
									numVideos,
									sourcesClient,
									log,
								)
								base.Progress(ctx, onProgress)
							}
						}
					}()
					if _, err := retrievePlaylistVideos(
						infoCache,
						e.ID,
						onVideo,
						entityLog,
					); err != nil {
						return errors.Wrap(err, "failed to retrieve playlist videos")
					}
				}
				base.Progress(ctx, onProgress)
				if err := cachePlaylistThumbnail(
					context.Background(), // "main" thread
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
					channel, err = infoCache.GetChannel(context.Background(), e.ID)
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
					go reportChannelDetails(channel, sourcesClient, log)
					onVideo := make(chan *api.VideoDetails, 4)
					go func() {
						numVideos := 0
						for {
							select {
							case <-done:
								return
							case video, ok := <-onVideo:
								if !ok {
									return
								}
								numVideos++
								go reportChannelVideo(
									e.ID,
									video,
									numVideos,
									sourcesClient,
									log,
								)
								base.Progress(ctx, onProgress)
							}
						}
					}()
					if _, err := retrieveChannelVideos(
						infoCache,
						e.ID,
						onVideo,
						entityLog,
					); err != nil {
						return errors.Wrap(err, "failed to retrieve channel videos")
					}
				}
				base.Progress(ctx, onProgress)
				if err := cacheChannelAvatar(
					context.Background(), // Background() since this is on the "main" thread
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
					go reportVideoDetails(video, sourcesClient, log)
				}
				base.Progress(ctx, onProgress)
				if err := cacheVideoThumbnail(
					context.Background(), // main thread
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
	details *api.VideoDetails,
	sourcesClient sources.Sources,
	log *zap.Logger,
) {
	if _, err := sourcesClient.ReportVideoDetails(
		context.TODO(),
		*convertVideo(details),
	); err != nil {
		log.Warn("failed to report video details", zap.Error(err))
	}
}

func reportPlaylistVideo(
	playlistID string,
	video *api.VideoDetails,
	numVideos int,
	sourcesClient sources.Sources,
	log *zap.Logger,
) {
	if _, err := sourcesClient.ReportPlaylistVideo(
		context.TODO(),
		sources.PlaylistVideo{
			PlaylistID: playlistID,
			Video:      *convertVideo(video),
			NumVideos:  numVideos,
		},
	); err != nil {
		log.Warn("failed to report playlist video", zap.Error(err))
	}
}

func reportChannelVideo(
	channelID string,
	video *api.VideoDetails,
	numVideos int,
	sourcesClient sources.Sources,
	log *zap.Logger,
) {
	if _, err := sourcesClient.ReportChannelVideo(
		context.TODO(),
		sources.ChannelVideo{
			ChannelID: channelID,
			NumVideos: numVideos,
			Video:     *convertVideo(video),
		},
	); err != nil {
		log.Warn("failed to report channel video", zap.Error(err))
	}
}

func reportChannelDetails(
	details *api.ChannelDetails,
	sourcesClient sources.Sources,
	log *zap.Logger,
) {
	if _, err := sourcesClient.ReportChannelDetails(
		context.TODO(),
		sources.ChannelDetails{
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
	details *api.PlaylistDetails,
	sourcesClient sources.Sources,
	log *zap.Logger,
) {
	if _, err := sourcesClient.ReportPlaylistDetails(
		context.TODO(),
		sources.PlaylistDetails{
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

func convertVideo(details *api.VideoDetails) *sources.VideoDetails {
	return &sources.VideoDetails{
		ID:          details.ID,
		Title:       details.Title,
		Description: details.Description,
		Thumbnail:   details.Thumbnail,
		UploadDate:  details.UploadDate,
		Uploader:    details.Uploader,
		UploaderID:  details.UploaderID,
		Channel:     details.Channel,
		ChannelID:   details.ChannelID,
		Duration:    details.Duration,
		ViewCount:   details.ViewCount,
		Width:       details.Width,
		Height:      details.Height,
		FPS:         details.FPS,
	}
}
