package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/thavlik/bjjvb/base/pkg/base"

	compiler "github.com/thavlik/bjjvb/compiler/pkg/api"

	"go.uber.org/zap"

	"github.com/pacedotdev/oto/otohttp"
	"github.com/thavlik/bjjvb/filter/pkg/api"
	"github.com/thavlik/bjjvb/filter/pkg/labelstore"
)

type Server struct {
	labelStore labelstore.LabelStore
	compiler   compiler.Compiler
	slideShow  base.ServiceOptions
	stackSize  int
	log        *zap.Logger
}

func NewServer(
	labelStore labelstore.LabelStore,
	compiler compiler.Compiler,
	slideShow base.ServiceOptions,
	stackSize int,
	log *zap.Logger,
) *Server {
	return &Server{
		labelStore,
		compiler,
		slideShow,
		stackSize,
		log,
	}
}

func (s *Server) ListenAndServe(port int) error {
	otoServer := otohttp.NewServer()
	api.RegisterFilter(otoServer, s)
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
