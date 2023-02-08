package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/pubsub"
	"github.com/thavlik/t4vd/base/pkg/scheduler"
	hound "github.com/thavlik/t4vd/hound/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/cachedset"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	"github.com/thavlik/t4vd/seer/pkg/thumbcache"
	"github.com/thavlik/t4vd/seer/pkg/vidcache"
	"github.com/thavlik/t4vd/seer/pkg/ytdl"
	"go.uber.org/zap"
)

var (
	cancelVideoTopic = "cancel_video"
)

func Entry(
	port int,
	pubSub pubsub.PubSub,
	querySched scheduler.Scheduler,
	dlSched scheduler.Scheduler,
	infoCache infocache.InfoCache,
	vidCache vidcache.VidCache,
	thumbCache thumbcache.ThumbCache,
	cachedVideoIDs cachedset.CachedSet,
	hound hound.Hound,
	dlOpts *ytdl.Options,
	concurrency int,
	disableDownloads bool,
	log *zap.Logger,
) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := NewServer(
		querySched,
		dlSched,
		pubSub,
		infoCache,
		vidCache,
		thumbCache,
		cachedVideoIDs,
		log,
	)

	stopPopQuery := make(chan struct{}, 1)
	defer func() { stopPopQuery <- struct{}{} }()
	initQueryWorkers(
		concurrency,
		infoCache,
		thumbCache,
		cachedVideoIDs,
		pubsub.Publisher(pubSub),
		querySched,
		hound,
		stopPopQuery,
		log,
	)

	cancelVideoDownload, err := pubSub.Subscribe(
		ctx,
		cancelVideoTopic,
	)
	if err != nil {
		return errors.Wrap(err, "failed to subscribe to topic")
	}

	stopPopDl := make(chan struct{}, 1)
	defer func() { stopPopDl <- struct{}{} }()
	initDownloadWorkers(
		concurrency,
		dlSched,
		cancelVideoDownload.Messages(ctx),
		vidCache,
		thumbCache,
		hound,
		dlOpts,
		disableDownloads,
		stopPopDl,
		log,
	)

	base.SignalReady(log)
	return s.ListenAndServe(port)
}
