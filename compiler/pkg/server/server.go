package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pacedotdev/oto/otohttp"
	"github.com/thavlik/bjjvb/base/pkg/base"
	"github.com/thavlik/bjjvb/base/pkg/scheduler"
	"github.com/thavlik/bjjvb/compiler/pkg/api"
	"github.com/thavlik/bjjvb/compiler/pkg/datastore"
	slideshow "github.com/thavlik/bjjvb/slideshow/pkg/api"
	sources "github.com/thavlik/bjjvb/sources/pkg/api"
	"go.uber.org/zap"
)

type Server struct {
	ds           datastore.DataStore
	sched        scheduler.Scheduler
	sources      sources.Sources
	seer         base.ServiceOptions
	slideshow    slideshow.SlideShow
	saveInterval time.Duration
	log          *zap.Logger
}

func NewServer(
	ds datastore.DataStore,
	scheduler scheduler.Scheduler,
	sources sources.Sources,
	seer base.ServiceOptions,
	slideshow slideshow.SlideShow,
	saveInterval time.Duration,
	log *zap.Logger,
) *Server {
	return &Server{
		ds:           ds,
		sched:        scheduler,
		sources:      sources,
		seer:         seer,
		slideshow:    slideshow,
		saveInterval: saveInterval,
		log:          log,
	}
}

func (s *Server) ListenAndServe(port int) error {
	otoServer := otohttp.NewServer()
	api.RegisterCompiler(otoServer, s)
	mux := http.NewServeMux()
	mux.Handle("/", otoServer)
	mux.HandleFunc("/readyz", base.ReadyHandler)
	s.log.Info("listening forever", zap.Int("port", port))
	return (&http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}).ListenAndServe()
}
