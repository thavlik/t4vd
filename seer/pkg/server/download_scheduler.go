package server

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/scheduler"
	hound "github.com/thavlik/t4vd/hound/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/thumbcache"
	"github.com/thavlik/t4vd/seer/pkg/vidcache"
	"go.uber.org/zap"
)

func initDownloadWorkers(
	concurrency int,
	dlSched scheduler.Scheduler,
	cancelVideoDownload <-chan []byte,
	vidCache vidcache.VidCache,
	thumbCache thumbcache.ThumbCache,
	hound hound.Hound,
	videoFormat string,
	includeAudio bool,
	disableDownloads bool,
	stop <-chan struct{},
	log *zap.Logger,
) {
	popVideoID := make(chan string)
	go downloadPopper(popVideoID, dlSched, stop, log)
	cancels := make([]chan []byte, concurrency)
	for i := 0; i < concurrency; i++ {
		cancel := make(chan []byte, 8)
		cancels[i] = cancel
		go downloadWorker(
			popVideoID,
			cancel,
			dlSched,
			vidCache,
			thumbCache,
			hound,
			videoFormat,
			includeAudio,
			disableDownloads,
			log,
		)
	}
	go func() {
		for {
			videoID, ok := <-cancelVideoDownload
			if !ok {
				return
			}
			for _, cancel := range cancels {
				cancel <- videoID
			}
		}
	}()
}

func downloadPopper(
	popVideoID chan<- string,
	dlSched scheduler.Scheduler,
	stop <-chan struct{},
	log *zap.Logger,
) {
	notification := dlSched.Notify()
	defer close(popVideoID)
	delay := 12 * time.Second
	for {
		start := time.Now()
		videoIDs, err := dlSched.List()
		if err != nil {
			panic(errors.Wrap(err, "scheduler.List"))
		}
		if len(videoIDs) > 0 {
			log.Debug("checking videos", zap.Int("num", len(videoIDs)))
			for _, videoID := range videoIDs {
				popVideoID <- videoID
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

func downloadWorker(
	popVideoID <-chan string,
	cancelVideoDownload <-chan []byte,
	dlSched scheduler.Scheduler,
	vidCache vidcache.VidCache,
	thumbCache thumbcache.ThumbCache,
	houndClient hound.Hound,
	videoFormat string,
	includeAudio bool,
	disableDownloads bool,
	log *zap.Logger,
) {
	for {
		videoID, ok := <-popVideoID
		if !ok {
			return
		}
		videoLog := log.With(zap.String("videoID", videoID))
		lock, err := dlSched.Lock(videoID)
		if err == scheduler.ErrLocked {
			// go to the next project
			videoLog.Debug("video already locked")
			continue
		} else if err != nil {
			panic(errors.Wrap(err, "scheduler.Lock"))
		}
		videoLog.Debug("locked video")
		func() {
			ctx, cancel := context.WithCancel(context.Background())
			stop := make(chan struct{}, 1)
			stopped := make(chan struct{})
			defer func() {
				stop <- struct{}{}
				<-stopped
				_ = lock.Release()
				cancel()
			}()
			onProgress := make(chan *base.DownloadProgress, 1)
			defer close(onProgress)
			go func() {
				defer func() { stopped <- struct{}{} }()
				done := ctx.Done()
				for {
					select {
					case cancelID := <-cancelVideoDownload:
						if string(cancelID) == videoID {
							cancel()
							videoLog.Debug("download was intentionally cancelled prematurely")
							return
						}
					case <-stop:
						return
					case <-done:
						return
					case progress, ok := <-onProgress:
						if !ok {
							return
						}
						if err := lock.Extend(); err != nil {
							videoLog.Warn("failed to extend video download lock", zap.Error(err))
							cancel()
							return
						}
						if progress != nil {
							go reportDownloadProgress(
								ctx,
								videoID,
								progress,
								houndClient,
								videoLog,
							)
						}
					}
				}
			}()
			if err := cacheVideoThumbnail(
				ctx,
				videoID,
				thumbCache,
				videoLog,
			); err != nil {
				videoLog.Error("error downloading thumbnail", zap.Error(err))
				return
			}
			base.ProgressDownload(ctx, onProgress)
			if err := downloadVideo(
				ctx,
				videoID,
				vidCache,
				videoFormat,
				includeAudio,
				disableDownloads,
				onProgress,
				videoLog,
			); err != nil {
				videoLog.Error("error downloading video", zap.Error(err))
				return
			}
			if err := dlSched.Remove(videoID); err != nil {
				videoLog.Warn("failed to remove videoID from download scheduler, this will result in multiple repeated requests to youtube")
			}
		}()
	}
}

func reportDownloadProgress(
	ctx context.Context,
	videoID string,
	progress *base.DownloadProgress,
	houndClient hound.Hound,
	log *zap.Logger,
) {
	if houndClient == nil {
		return
	}
	if _, err := houndClient.ReportVideoDownloadProgress(
		ctx,
		hound.VideoDownloadProgress{
			ID:      videoID,
			Total:   progress.Total,
			Rate:    progress.Rate,
			Elapsed: int64(progress.Elapsed),
		},
	); err != nil {
		log.Warn("failed to report download progress", zap.Error(err))
	}
}
