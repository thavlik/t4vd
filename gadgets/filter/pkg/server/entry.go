package server

import (
	"github.com/thavlik/t4vd/filter/pkg/labelstore"

	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/gadget"
	"go.uber.org/zap"
)

func Entry(
	port int,
	labelStore labelstore.LabelStore,
	defaultRef *gadget.DataRef,
	maxBatchSize int,
	log *zap.Logger,
) error {
	base.SignalReady(log)
	return NewServer(
		labelStore,
		maxBatchSize,
		defaultRef,
		log,
	).Listen(port)
}
