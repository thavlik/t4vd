package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/thavlik/t4vd/base/pkg/base"
	compiler "github.com/thavlik/t4vd/compiler/pkg/api"

	"go.uber.org/zap"

	"github.com/pacedotdev/oto/otohttp"
	"github.com/thavlik/t4vd/slideshow/pkg/api"
	"github.com/thavlik/t4vd/slideshow/pkg/imgcache"
	"github.com/thavlik/t4vd/slideshow/pkg/markercache"
)

type Server struct {
	bucket      string
	imgCache    imgcache.ImgCache
	markerCache markercache.MarkerCache
	compiler    compiler.Compiler
	log         *zap.Logger
}

func NewServer(
	bucket string,
	imgCache imgcache.ImgCache,
	markerCache markercache.MarkerCache,
	compiler compiler.Compiler,
	log *zap.Logger,
) *Server {
	return &Server{
		bucket:      bucket,
		imgCache:    imgCache,
		markerCache: markerCache,
		log:         log,
	}
}

func (s *Server) ListenAndServe(port int) error {
	otoServer := otohttp.NewServer()
	api.RegisterSlideShow(otoServer, s)
	mux := http.NewServeMux()
	mux.Handle("/", otoServer)
	mux.HandleFunc("/healthz", base.HealthHandler)
	mux.HandleFunc("/readyz", base.ReadyHandler)
	mux.HandleFunc("/frame", s.handleGetFrame())
	s.log.Info("listening forever", zap.Int("port", port))
	return (&http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}).ListenAndServe()
}
