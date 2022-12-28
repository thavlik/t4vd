package server

import (
	"github.com/thavlik/bjjvb/base/pkg/base"
	compiler "github.com/thavlik/bjjvb/compiler/pkg/api"
	"github.com/thavlik/bjjvb/filter/pkg/labelstore"
	"go.uber.org/zap"
)

func Entry(
	port int,
	labelStore labelstore.LabelStore,
	compiler compiler.Compiler,
	slideshow base.ServiceOptions,
	stackSize int,
	log *zap.Logger,
) error {
	s := NewServer(
		labelStore,
		compiler,
		slideshow,
		stackSize,
		log,
	)
	base.SignalReady(log)
	err := s.ListenAndServe(port)
	return err
}
