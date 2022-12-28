package server

import (
	"github.com/thavlik/bjjvb/base/pkg/base"
	"github.com/thavlik/bjjvb/base/pkg/iam"
	compiler "github.com/thavlik/bjjvb/compiler/pkg/api"
	filter "github.com/thavlik/bjjvb/filter/pkg/api"
	sources "github.com/thavlik/bjjvb/sources/pkg/api"
	"go.uber.org/zap"
)

func Entry(
	port int,
	adminPort int,
	iam iam.IAM,
	seerOpts base.ServiceOptions,
	sources sources.Sources,
	compiler compiler.Compiler,
	filter filter.Filter,
	slideshow base.ServiceOptions,
	corsHeader string,
	log *zap.Logger,
) error {
	base.SignalReady(log)
	s := NewServer(
		iam,
		seerOpts,
		sources,
		compiler,
		filter,
		slideshow,
		corsHeader,
		log,
	)
	mainErr := make(chan error)
	go func() {
		mainErr <- s.ListenAndServe(port)
	}()
	adminErr := make(chan error)
	go func() {
		mainErr <- s.AdminListenAndServe(adminPort)
	}()
	select {
	case err := <-mainErr:
		return err
	case err := <-adminErr:
		return err
	}
}
