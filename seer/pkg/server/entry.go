package server

import (
	"github.com/thavlik/bjjvb/base/pkg/base"
	"github.com/thavlik/bjjvb/base/pkg/scheduler"
	"github.com/thavlik/bjjvb/seer/pkg/infocache"
	"github.com/thavlik/bjjvb/seer/pkg/thumbcache"
	"github.com/thavlik/bjjvb/seer/pkg/vidcache"
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
		stopPopDl,
		log,
	)

	base.SignalReady(log)
	return s.ListenAndServe(port)
}
