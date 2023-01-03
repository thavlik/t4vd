package server

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/seer/pkg/vidcache"
	"github.com/thavlik/t4vd/seer/pkg/ytdl"
	"go.uber.org/zap"
)

func downloadVideo(
	ctx context.Context,
	videoID string,
	vidCache vidcache.VidCache,
	videoFormat string,
	includeAudio bool,
	disableDownloads bool,
	onProgress chan<- *base.DownloadProgress,
	log *zap.Logger,
) error {
	// we can get away with only checking the
	// object's head to see if it exists
	if has, err := vidCache.Has(ctx, videoID); err != nil {
		return errors.Wrap(err, "cache.Has")
	} else if has {
		// cache already has the video and we don't
		// want to download it here
		log.Debug("cache has video")
		return nil
	}
	if disableDownloads {
		log.Warn("video would be downloaded from youtube but --disable-downloads was specified")
		return nil
	}
	log.Debug("downloading video from youtube")
	rp, wp := io.Pipe()
	ytdlDone := make(chan error, 1)
	// we want the download to terminate when it's
	// complete but we don't want to terminate if
	// the request times out
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		ytdlDone <- ytdl.Download(
			ctx,
			videoID,
			wp,
			videoFormat,
			includeAudio,
			onProgress,
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
		return errors.Wrap(err, "ytdl.Download")
	}
	log.Debug("cached video")
	base.ProgressDownload(ctx, onProgress)
	return nil
}
