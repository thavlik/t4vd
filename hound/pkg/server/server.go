package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/thavlik/t4vd/base/pkg/base"

	compiler "github.com/thavlik/t4vd/compiler/pkg/api"
	gateway "github.com/thavlik/t4vd/gateway/pkg/api"
	sources "github.com/thavlik/t4vd/sources/pkg/api"

	"go.uber.org/zap"

	"github.com/pacedotdev/oto/otohttp"
	"github.com/thavlik/t4vd/hound/pkg/api"
)

type Server struct {
	compiler compiler.Compiler
	sources  sources.Sources
	gateway  gateway.Gateway
	log      *zap.Logger
}

func NewServer(
	compiler compiler.Compiler,
	sources sources.Sources,
	gateway gateway.Gateway,
	log *zap.Logger,
) *Server {
	return &Server{
		compiler,
		sources,
		gateway,
		log,
	}
}

func (s *Server) ListenAndServe(port int) error {
	otoServer := otohttp.NewServer()
	api.RegisterHound(otoServer, s)
	mux := http.NewServeMux()
	mux.Handle("/", otoServer)
	mux.HandleFunc("/healthz", base.HealthHandler)
	mux.HandleFunc("/readyz", base.ReadyHandler)
	s.log.Info("listening forever", zap.Int("port", port))
	return (&http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}).ListenAndServe()
}
