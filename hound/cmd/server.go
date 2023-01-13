package main

import (
	"time"

	"github.com/thavlik/t4vd/base/pkg/base"

	"github.com/spf13/cobra"
	compiler "github.com/thavlik/t4vd/compiler/pkg/api"
	gateway "github.com/thavlik/t4vd/gateway/pkg/api"
	"github.com/thavlik/t4vd/hound/pkg/server"
	sources "github.com/thavlik/t4vd/sources/pkg/api"
)

var defaultTimeout = 8 * time.Second

var serverArgs struct {
	base.ServerOptions
	compiler base.ServiceOptions
	gateway  base.ServiceOptions
	sources  base.ServiceOptions
}

var serverCmd = &cobra.Command{
	Use: "server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.ServerEnv(&serverArgs.ServerOptions)
		base.ServiceEnv("compiler", &serverArgs.compiler)
		base.ServiceEnv("gateway", &serverArgs.gateway)
		base.ServiceEnv("sources", &serverArgs.sources)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log := base.DefaultLog
		go base.RunMetrics(serverArgs.MetricsPort, log)
		base.RandomizeSeed()
		return server.Entry(
			serverArgs.Port,
			compiler.NewCompilerClientFromOptions(serverArgs.compiler),
			sources.NewSourcesClientFromOptions(serverArgs.sources),
			gateway.NewGatewayClientFromOptions(serverArgs.gateway),
			log,
		)
	},
}

func init() {
	base.AddServerFlags(serverCmd, &serverArgs.ServerOptions)
	base.AddServiceFlags(serverCmd, "compiler", &serverArgs.compiler, defaultTimeout)
	base.AddServiceFlags(serverCmd, "sources", &serverArgs.sources, defaultTimeout)
	base.AddServiceFlags(serverCmd, "gateway", &serverArgs.gateway, defaultTimeout)
	ConfigureCommand(serverCmd)
}
