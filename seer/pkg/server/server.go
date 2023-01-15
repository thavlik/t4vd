package server

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/pacedotdev/oto/otohttp"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/pubsub"
	"github.com/thavlik/t4vd/base/pkg/scheduler"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/cachedset"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	"github.com/thavlik/t4vd/seer/pkg/thumbcache"
	"github.com/thavlik/t4vd/seer/pkg/vidcache"
)

type Server struct {
	querySched     scheduler.Scheduler
	dlSched        scheduler.Scheduler
	pubsub         pubsub.PubSub
	infoCache      infocache.InfoCache
	vidCache       vidcache.VidCache
	thumbCache     thumbcache.ThumbCache
	cachedVideoIDs cachedset.CachedSet
	videoFormat    string
	includeAudio   bool
	log            *zap.Logger
}

func NewServer(
	querySched scheduler.Scheduler,
	dlSched scheduler.Scheduler,
	pub pubsub.PubSub,
	infoCache infocache.InfoCache,
	vidCache vidcache.VidCache,
	thumbCache thumbcache.ThumbCache,
	cachedVideoIDs cachedset.CachedSet,
	videoFormat string,
	includeAudio bool,
	log *zap.Logger,
) *Server {
	return &Server{
		querySched,
		dlSched,
		pub,
		infoCache,
		vidCache,
		thumbCache,
		cachedVideoIDs,
		videoFormat,
		includeAudio,
		log,
	}
}

func (s *Server) ListenAndServe(port int) error {
	otoServer := otohttp.NewServer()
	api.RegisterSeer(otoServer, s)
	mux := http.NewServeMux()
	mux.Handle("/", otoServer)
	mux.HandleFunc("/channel/videos", s.handleGetChannelVideoIDs())
	mux.HandleFunc("/playlist/videos", s.handleGetPlaylistVideoIDs())
	mux.HandleFunc("/video", s.handleGetVideo())
	mux.HandleFunc("/video/thumbnail", s.handleGetVideoThumbnail())
	mux.HandleFunc("/playlist/thumbnail", s.handleGetPlaylistThumbnail())
	mux.HandleFunc("/channel/avatar", s.handleGetChannelAvatar())
	mux.HandleFunc("/healthz", base.HealthHandler)
	mux.HandleFunc("/readyz", base.ReadyHandler)
	s.log.Info("listening forever", zap.Int("port", port))
	return (&http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: time.Hour,
	}).ListenAndServe()
}
