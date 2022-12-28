package server

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/base/pkg/base"
	"github.com/thavlik/bjjvb/seer/pkg/vidcache"
	"github.com/thavlik/bjjvb/seer/pkg/ytdl"
	"go.uber.org/zap"
)

func downloadVideo(
	ctx context.Context,
	videoID string,
	w io.Writer,
	vidCache vidcache.VidCache,
	videoFormat string,
	includeAudio bool,
	onProgress chan<- struct{},
	log *zap.Logger,
) error {
	defer base.Progress(ctx, onProgress)
	noDownload := w == nil
	if noDownload {
		// we can get away with only checking the
		// object's head to see if it exists
		if has, err := vidCache.Has(ctx, videoID); err != nil {
			return errors.Wrap(err, "cache.Has")
		} else if has {
			// cache already has the video and we don't
			// want to download it here
			log.Debug("cache has video")
			return nil
		} else {
			log.Debug("video not cached")
		}
	} else {
		// try and get the cached video
		if err := vidCache.Get(ctx, videoID, w); err == nil {
			log.Debug("served video from cache")
			return nil
		} else if err != vidcache.ErrVideoNotCached {
			return errors.Wrap(err, "cache.Get")
		}
	}
	log.Debug("downloading video from youtube")
	rp, wp := io.Pipe()
	ytdlDone := make(chan error, 1)
	// we want the download to terminate when it's
	// complete but we don't want to terminate if
	// the request times out
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		var onProg chan struct{}
		if onProgress != nil {
			stop := make(chan struct{}, 1)
			stopped := make(chan struct{})
			defer func() {
				stop <- struct{}{}
				<-stopped
			}()
			onProg = make(chan struct{})
			go func() {
				defer func() { stopped <- struct{}{} }()
				for {
					select {
					case <-stop:
						return
					case _, ok := <-onProg:
						if !ok {
							return
						}
						base.Progress(ctx, onProgress)
					}
				}
			}()
		}
		ytdlDone <- ytdl.Download(
			ctx,
			videoID,
			wp,
			videoFormat,
			includeAudio,
			onProg,
			log,
		)
		close(ytdlDone)
	}()
	err := vidCache.Set(videoID, rp)
	_ = rp.Close()
	_ = wp.Close()
	if err != nil {
		return errors.Wrap(err, "set cache")
	}
	log.Debug("waiting on ytdl")
	err = <-ytdlDone
	log.Debug("finished video download")
	if err != nil {
		return errors.Wrap(err, "ytdl")
	}
	log.Debug("cached video")
	base.Progress(ctx, onProgress)
	if !noDownload {
		// finally download the video from the cache
		if err := vidCache.Get(ctx, videoID, w); err != nil {
			return errors.Wrap(err, "get cache")
		}
	}
	return nil
}
