package server

import (
	"github.com/thavlik/t4vd/base/pkg/base"
	compiler "github.com/thavlik/t4vd/compiler/pkg/api"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
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
