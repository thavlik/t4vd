package main

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/thavlik/bjjvb/base/cmd/iam"
	"github.com/thavlik/bjjvb/base/pkg/base"
	compiler "github.com/thavlik/bjjvb/compiler/pkg/api"
	filter "github.com/thavlik/bjjvb/filter/pkg/api"
	"github.com/thavlik/bjjvb/gateway/pkg/server"
	sources "github.com/thavlik/bjjvb/sources/pkg/api"
)

var defaultTimeout = 10 * time.Second

var serverArgs struct {
	base.ServerOptions
	adminPort  int
	sources    base.ServiceOptions
	iam        base.IAMOptions
	compiler   base.ServiceOptions
	filter     base.ServiceOptions
	slideshow  base.ServiceOptions
	seer       base.ServiceOptions
	corsHeader string
}

var serverCmd = &cobra.Command{
	Use: "server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.ServerEnv(&serverArgs.ServerOptions)
		base.ServiceEnv("sources", &serverArgs.sources)
		base.ServiceEnv("compiler", &serverArgs.compiler)
		base.ServiceEnv("filter", &serverArgs.filter)
		base.ServiceEnv("slide-show", &serverArgs.slideshow)
		base.ServiceEnv("seer", &serverArgs.seer)
		base.IAMEnv(&serverArgs.iam, false)
		base.CheckEnv("CORS_HEADER", &serverArgs.corsHeader)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log := base.Log
		go base.RunMetrics(serverArgs.MetricsPort, log)
		return server.Entry(
			serverArgs.Port,
			serverArgs.adminPort,
			iam.InitIAM(&serverArgs.iam),
			serverArgs.seer,
			sources.NewSourcesClientFromOptions(serverArgs.sources),
			compiler.NewCompilerClientFromOptions(serverArgs.compiler),
			filter.NewFilterClientFromOptions(serverArgs.filter),
			serverArgs.slideshow,
			serverArgs.corsHeader,
			log,
		)
	},
}

func init() {
	base.AddServerFlags(serverCmd, &serverArgs.ServerOptions)
	serverCmd.PersistentFlags().IntVar(&serverArgs.adminPort, "admin-port", 8080, "http service port")
	serverCmd.PersistentFlags().StringVar(&serverArgs.corsHeader, "cors-header", "", "Access-Control-Allow-Origin header")
	base.AddServiceFlags(serverCmd, "compiler", &serverArgs.compiler, defaultTimeout)
	base.AddServiceFlags(serverCmd, "sources", &serverArgs.sources, defaultTimeout)
	base.AddServiceFlags(serverCmd, "filter", &serverArgs.filter, defaultTimeout)
	base.AddServiceFlags(serverCmd, "slide-show", &serverArgs.slideshow, defaultTimeout)
	base.AddServiceFlags(serverCmd, "seer", &serverArgs.seer, defaultTimeout)
	base.AddIAMFlags(serverCmd, &serverArgs.iam)
	ConfigureCommand(serverCmd)
}
