package server

import (
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/iam"
	compiler "github.com/thavlik/t4vd/compiler/pkg/api"
	seer "github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/sources/pkg/store"
	"go.uber.org/zap"
)

func Entry(
	port int,
	iam iam.IAM,
	store store.Store,
	seer seer.Seer,
	compiler compiler.Compiler,
	log *zap.Logger,
) error {
	base.SignalReady(log)
	return NewServer(
		iam,
		store,
		seer,
		compiler,
		log,
	).ListenAndServe(port)
}
