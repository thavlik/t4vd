package main

import (
	"time"

	"github.com/thavlik/t4vd/base/pkg/base"

	"github.com/spf13/cobra"
	compiler "github.com/thavlik/t4vd/compiler/pkg/api"
	gateway "github.com/thavlik/t4vd/gateway/pkg/api"
	"github.com/thavlik/t4vd/hound/pkg/server"
)

var serverArgs struct {
	base.ServerOptions
	compiler base.ServiceOptions
	gateway  base.ServiceOptions
}

var serverCmd = &cobra.Command{
	Use: "server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.ServerEnv(&serverArgs.ServerOptions)
		base.ServiceEnv("compiler", &serverArgs.compiler)
		base.ServiceEnv("gateway", &serverArgs.gateway)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		go base.RunMetrics(serverArgs.MetricsPort, base.Log)
		base.RandomizeSeed()
		return server.Entry(
			serverArgs.Port,
			compiler.NewCompilerClientFromOptions(serverArgs.compiler),
			gateway.NewGatewayClientFromOptions(serverArgs.gateway),
			base.Log,
		)
	},
}

func init() {
	base.AddServerFlags(serverCmd, &serverArgs.ServerOptions)
	base.AddServiceFlags(serverCmd, "compiler", &serverArgs.compiler, 8*time.Second)
	base.AddServiceFlags(serverCmd, "gateway", &serverArgs.gateway, 8*time.Second)
	ConfigureCommand(serverCmd)
}
