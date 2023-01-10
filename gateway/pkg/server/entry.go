package server

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/iam"
	"github.com/thavlik/t4vd/base/pkg/pubsub"
	compiler "github.com/thavlik/t4vd/compiler/pkg/api"
	filter "github.com/thavlik/t4vd/filter/pkg/api"
	"github.com/thavlik/t4vd/gateway/pkg/api"
	sources "github.com/thavlik/t4vd/sources/pkg/api"
	"go.uber.org/zap"
)

var gatewayTopic = "gateway"

func Entry(
	port int,
	adminPort int,
	iam iam.IAM,
	seerOpts base.ServiceOptions,
	sources sources.Sources,
	compiler compiler.Compiler,
	filter filter.Filter,
	slideshow base.ServiceOptions,
	pubSub pubsub.PubSub,
	corsHeader string,
	log *zap.Logger,
) error {
	s := NewServer(
		iam,
		seerOpts,
		sources,
		compiler,
		filter,
		slideshow,
		pubsub.Publisher(pubSub),
		corsHeader,
		log,
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sub, err := pubSub.Subscribe(ctx, gatewayTopic)
	if err != nil {
		return errors.Wrap(err, "pubsub.Subscribe")
	}
	go func() {
		defer panic("gateway subscriber returned")
		ch := sub.Messages(ctx)
		defer sub.Cancel(ctx)
		done := ctx.Done()
		for {
			select {
			case <-done:
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				var event api.Event
				if err := json.Unmarshal(msg, &event); err != nil {
					panic(errors.Wrap(err, "json.Unmarshal"))
				}
				if err := s.pushEventLocal(event); err != nil {
					panic(errors.Wrap(err, "pushEventLocal"))
				}
			}
		}
	}()
	mainErr := make(chan error, 1)
	go func() { mainErr <- s.ListenAndServe(port) }()
	adminErr := make(chan error, 1)
	go func() { adminErr <- s.AdminListenAndServe(adminPort) }()
	base.SignalReady(log)
	select {
	case err := <-mainErr:
		return err
	case err := <-adminErr:
		return err
	}
}
