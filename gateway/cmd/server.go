package main

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/thavlik/t4vd/base/cmd/iam"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/pubsub"
	memory_pubsub "github.com/thavlik/t4vd/base/pkg/pubsub/memory"
	redis_pubsub "github.com/thavlik/t4vd/base/pkg/pubsub/redis"
	compiler "github.com/thavlik/t4vd/compiler/pkg/api"
	filter "github.com/thavlik/t4vd/filter/pkg/api"
	"github.com/thavlik/t4vd/gateway/pkg/server"
	sources "github.com/thavlik/t4vd/sources/pkg/api"
	"go.uber.org/zap"
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
	redis      base.RedisOptions
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
		base.RedisEnv(&serverArgs.redis, false)
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
			initPubSub(log),
			serverArgs.corsHeader,
			log,
		)
	},
}

func initPubSub(log *zap.Logger) pubsub.PubSub {
	if serverArgs.redis.IsSet() {
		return redis_pubsub.NewRedisPubSub(
			base.ConnectRedis(&serverArgs.redis),
			"gateway",
			log,
		)
	}
	return memory_pubsub.NewMemoryPubSub(log)
}

func init() {
	base.AddServerFlags(serverCmd, &serverArgs.ServerOptions)
	base.AddRedisFlags(serverCmd, &serverArgs.redis)
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
