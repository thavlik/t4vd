package server

import (
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/scheduler"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	"github.com/thavlik/t4vd/seer/pkg/thumbcache"
	"github.com/thavlik/t4vd/seer/pkg/vidcache"
	"go.uber.org/zap"
)

func Entry(
	port int,
	querySched scheduler.Scheduler,
	dlSched scheduler.Scheduler,
	infoCache infocache.InfoCache,
	vidCache vidcache.VidCache,
	thumbCache thumbcache.ThumbCache,
	videoFormat string,
	includeAudio bool,
	concurrency int,
	disableDownloads bool,
	log *zap.Logger,
) error {
	s := NewServer(
		querySched,
		dlSched,
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
		stopPopQuery,
		log,
	)

	stopPopDl := make(chan struct{}, 1)
	defer func() { stopPopDl <- struct{}{} }()
	initDownloadWorkers(
		concurrency,
		dlSched,
		vidCache,
		thumbCache,
		videoFormat,
		includeAudio,
		disableDownloads,
		stopPopDl,
		log,
	)

	base.SignalReady(log)
	return s.ListenAndServe(port)
}
