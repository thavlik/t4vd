package main

import (
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/thavlik/bjjvb/base/pkg/base"
	compiler "github.com/thavlik/bjjvb/compiler/pkg/api"
	"github.com/thavlik/bjjvb/slideshow/pkg/imgcache"
	memory_imgcache "github.com/thavlik/bjjvb/slideshow/pkg/imgcache/memory"
	redis_imgcache "github.com/thavlik/bjjvb/slideshow/pkg/imgcache/redis"
	"github.com/thavlik/bjjvb/slideshow/pkg/server"
	sources "github.com/thavlik/bjjvb/sources/pkg/api"
)

var serverArgs struct {
	base.ServerOptions
	compiler base.ServiceOptions
	sources  base.ServiceOptions
	seer     base.ServiceOptions
	redis    base.RedisOptions
	bucket   string
}

var serverCmd = &cobra.Command{
	Use: "server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.ServerEnv(&serverArgs.ServerOptions)
		base.ServiceEnv("compiler", &serverArgs.compiler)
		base.ServiceEnv("seer", &serverArgs.seer)
		base.ServiceEnv("sources", &serverArgs.sources)
		base.RedisEnv(&serverArgs.redis, false)
		if serverArgs.bucket == "" {
			return errors.New("missing --bucket")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log := base.Log
		go base.RunMetrics(serverArgs.MetricsPort, log)
		base.RandomizeSeed()
		var redis *redis.Client
		if serverArgs.redis.IsSet() {
			redis = base.ConnectRedis(&serverArgs.redis)
		}
		return server.Entry(
			serverArgs.Port,
			serverArgs.bucket,
			imgCacheClient(redis),
			redis,
			compiler.NewCompilerClientFromOptions(serverArgs.compiler),
			sources.NewSourcesClientFromOptions(serverArgs.sources),
			log,
		)
	},
}

func imgCacheClient(redis *redis.Client) imgcache.ImgCache {
	if redis == nil {
		return memory_imgcache.NewMemoryImgCache(
			base.CheckEnvInt("IMG_CACHE_SIZE", nil))
	}
	return redis_imgcache.NewRedisImgCache(redis, 0)
}

func init() {
	base.AddServerFlags(serverCmd, &serverArgs.ServerOptions)
	base.AddServiceFlags(serverCmd, "compiler", &serverArgs.compiler, time.Minute)
	base.AddServiceFlags(serverCmd, "sources", &serverArgs.sources, 12*time.Second)
	base.AddServiceFlags(serverCmd, "seer", &serverArgs.seer, 8*time.Second)
	base.AddRedisFlags(serverCmd, &serverArgs.redis)
	serverCmd.PersistentFlags().StringVar(&serverArgs.bucket, "bucket", "", "the bucket containing full length webm videos")
	//serverCmd.PersistentFlags().IntVar(&serverArgs.markerBufferSize, "marker-buffer-size", 0, "buffer for frame-cached markers")
	ConfigureCommand(serverCmd)
}
