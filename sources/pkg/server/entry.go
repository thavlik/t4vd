package server

import (
	"github.com/thavlik/bjjvb/base/pkg/base"
	"github.com/thavlik/bjjvb/base/pkg/iam"
	compiler "github.com/thavlik/bjjvb/compiler/pkg/api"
	seer "github.com/thavlik/bjjvb/seer/pkg/api"
	"github.com/thavlik/bjjvb/sources/pkg/store"
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
