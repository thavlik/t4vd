package server

import (
	"github.com/thavlik/t4vd/base/pkg/base"
	compiler "github.com/thavlik/t4vd/compiler/pkg/api"
	gateway "github.com/thavlik/t4vd/gateway/pkg/api"
	"go.uber.org/zap"
)

func Entry(
	port int,
	compiler compiler.Compiler,
	gateway gateway.Gateway,
	log *zap.Logger,
) error {
	s := NewServer(
		compiler,
		gateway,
		log,
	)
	base.SignalReady(log)
	return s.ListenAndServe(port)
}
