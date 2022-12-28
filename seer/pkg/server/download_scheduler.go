package server

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/scheduler"
	"github.com/thavlik/t4vd/seer/pkg/thumbcache"
	"github.com/thavlik/t4vd/seer/pkg/vidcache"
	"go.uber.org/zap"
)

func initDownloadWorkers(
	concurrency int,
	dlSched scheduler.Scheduler,
	vidCache vidcache.VidCache,
	thumbCache thumbcache.ThumbCache,
	videoFormat string,
	includeAudio bool,
	stop <-chan struct{},
	log *zap.Logger,
) {
	popVideoID := make(chan string)
	go downloadPopper(popVideoID, dlSched, stop, log)
	for i := 0; i < concurrency; i++ {
		go downloadWorker(
			popVideoID,
			dlSched,
			vidCache,
			thumbCache,
			videoFormat,
			includeAudio,
			log,
		)
	}
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
	dlSched scheduler.Scheduler,
	vidCache vidcache.VidCache,
	thumbCache thumbcache.ThumbCache,
	videoFormat string,
	includeAudio bool,
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
			onProgress := make(chan struct{}, 1)
			go func() {
				defer func() { stopped <- struct{}{} }()
				done := ctx.Done()
				for {
					select {
					case <-stop:
						return
					case <-done:
						return
					case _, ok := <-onProgress:
						if !ok {
							return
						}
						if err := lock.Extend(); err != nil {
							videoLog.Warn("failed to extend video download lock", zap.Error(err))
							cancel()
							return
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
			base.Progress(ctx, onProgress)
			if err := downloadVideo(
				ctx,
				videoID,
				nil,
				vidCache,
				videoFormat,
				includeAudio,
				onProgress,
				videoLog,
			); err != nil {
				videoLog.Error("error downloading video", zap.Error(err))
				return
			}
			base.Progress(ctx, onProgress)
			if err := dlSched.Remove(videoID); err != nil {
				videoLog.Warn("failed to remove videoID from download scheduler, this will result in multiple repeated requests to youtube")
			}
		}()
	}
}
