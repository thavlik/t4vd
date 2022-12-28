package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/thavlik/bjjvb/base/pkg/base"
	"github.com/thavlik/bjjvb/base/pkg/iam"
	compiler "github.com/thavlik/bjjvb/compiler/pkg/api"

	seer "github.com/thavlik/bjjvb/seer/pkg/api"

	"go.uber.org/zap"

	"github.com/pacedotdev/oto/otohttp"
	"github.com/thavlik/bjjvb/sources/pkg/api"
	"github.com/thavlik/bjjvb/sources/pkg/store"
)

var (
	errMissingProjectID   = errors.New("missing project id")
	errMissingSubmitterID = errors.New("missing submitter id")
)

type Server struct {
	iam      iam.IAM
	store    store.Store
	seer     seer.Seer
	compiler compiler.Compiler
	log      *zap.Logger
}

func NewServer(
	iam iam.IAM,
	store store.Store,
	seer seer.Seer,
	compiler compiler.Compiler,
	log *zap.Logger,
) *Server {
	return &Server{
		iam,
		store,
		seer,
		compiler,
		log,
	}
}

func (s *Server) ListenAndServe(port int) error {
	otoServer := otohttp.NewServer()
	api.RegisterSources(otoServer, s)
	s.log.Info("listening forever", zap.Int("port", port))
	mux := http.NewServeMux()
	mux.Handle("/", otoServer)
	mux.HandleFunc("/readyz", base.ReadyHandler)
	return (&http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}).ListenAndServe()
}
