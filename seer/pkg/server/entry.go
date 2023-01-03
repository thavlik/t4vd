package server

import (
	"context"

	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/pubsub"
	"github.com/thavlik/t4vd/base/pkg/scheduler"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	"github.com/thavlik/t4vd/seer/pkg/thumbcache"
	"github.com/thavlik/t4vd/seer/pkg/vidcache"
	sources "github.com/thavlik/t4vd/sources/pkg/api"
	"go.uber.org/zap"
)

func Entry(
	port int,
	pubSub pubsub.PubSub,
	querySched scheduler.Scheduler,
	dlSched scheduler.Scheduler,
	infoCache infocache.InfoCache,
	vidCache vidcache.VidCache,
	thumbCache thumbcache.ThumbCache,
	sources sources.Sources,
	videoFormat string,
	includeAudio bool,
	concurrency int,
	disableDownloads bool,
	log *zap.Logger,
) error {
	s := NewServer(
		querySched,
		dlSched,
		pubsub.Publisher(pubSub),
		infoCache,
		vidCache,
		thumbCache,
		videoFormat,
		includeAudio,
		log,
	)

	stopPopQuery := make(chan struct{}, 1)
	defer func() { stopPopQuery <- struct{}{} }()
	initQueryWorkers(
		concurrency,
		infoCache,
		thumbCache,
		querySched,
		sources,
		stopPopQuery,
		log,
	)

	cancelVideoDownload, err := pubSub.Subscribe(context.Background())
	if err != nil {
		return err
	}

	stopPopDl := make(chan struct{}, 1)
	defer func() { stopPopDl <- struct{}{} }()
	initDownloadWorkers(
		concurrency,
		dlSched,
		cancelVideoDownload,
		vidCache,
		thumbCache,
		sources,
		videoFormat,
		includeAudio,
		disableDownloads,
		stopPopDl,
		log,
	)

	base.SignalReady(log)
	return s.ListenAndServe(port)
}
