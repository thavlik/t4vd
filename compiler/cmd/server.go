package main

import (
	"time"

	"github.com/go-redis/redis/v9"
	slideshow "github.com/thavlik/t4vd/slideshow/pkg/api"
	sources "github.com/thavlik/t4vd/sources/pkg/api"
	"go.uber.org/zap"

	"github.com/thavlik/t4vd/base/pkg/base"

	"github.com/spf13/cobra"
	"github.com/thavlik/t4vd/base/pkg/scheduler"
	memory_scheduler "github.com/thavlik/t4vd/base/pkg/scheduler/memory"
	redis_scheduler "github.com/thavlik/t4vd/base/pkg/scheduler/redis"
	"github.com/thavlik/t4vd/compiler/pkg/datacache"
	redis_datacache "github.com/thavlik/t4vd/compiler/pkg/datacache/redis"
	"github.com/thavlik/t4vd/compiler/pkg/datastore"
	mongo_datastore "github.com/thavlik/t4vd/compiler/pkg/datastore/mongo"
	postgres_datastore "github.com/thavlik/t4vd/compiler/pkg/datastore/postgres"
	"github.com/thavlik/t4vd/compiler/pkg/server"
)

var serverArgs struct {
	base.ServerOptions
	redis            base.RedisOptions
	sources          base.ServiceOptions
	seer             base.ServiceOptions
	slideshow        base.ServiceOptions
	db               base.DatabaseOptions
	saveInterval     time.Duration
	concurrency      int
	compileOnStart   bool
	maxVideoDuration time.Duration
}

var serverCmd = &cobra.Command{
	Use: "server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.ServerEnv(&serverArgs.ServerOptions)
		base.RedisEnv(&serverArgs.redis, false)
		base.ServiceEnv("seer", &serverArgs.seer)
		base.ServiceEnv("sources", &serverArgs.sources)
		base.ServiceEnv("slide-show", &serverArgs.slideshow)
		base.DatabaseEnv(&serverArgs.db, true)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log := base.DefaultLog
		var redis *redis.Client
		if serverArgs.redis.IsSet() {
			redis = base.ConnectRedis(&serverArgs.redis)
		}
		go base.RunMetrics(serverArgs.MetricsPort, log)
		return server.Entry(
			serverArgs.Port,
			initScheduler(redis),
			initDataStore(),
			initDataCache(redis, log),
			sources.NewSourcesClientFromOptions(serverArgs.sources),
			serverArgs.seer,
			slideshow.NewSlideShowClientFromOptions(serverArgs.slideshow),
			serverArgs.saveInterval,
			serverArgs.compileOnStart,
			serverArgs.concurrency,
			log,
		)
	},
}

func initDataCache(
	redis *redis.Client,
	log *zap.Logger,
) datacache.DataCache {
	if redis != nil {
		return redis_datacache.NewRedisDataCache(
			redis,
			log,
		)
	}
	panic("only redis datacache is currently supported")
}

func initScheduler(redis *redis.Client) scheduler.Scheduler {
	if redis != nil {
		return redis_scheduler.NewRedisScheduler(
			redis,
			"compsched",
			10*time.Second,
		)
	}
	return memory_scheduler.NewMemoryScheduler()
}

func initDataStore() datastore.DataStore {
	log := base.DefaultLog
	switch serverArgs.db.Driver {
	case base.PostgresDriver:
		ds, err := postgres_datastore.NewPostgresDataStore(
			base.ConnectPostgres(&serverArgs.db.Postgres),
			log,
		)
		if err != nil {
			panic(err)
		}
		return ds
	case base.MongoDriver:
		return mongo_datastore.NewMongoDataStore(
			base.ConnectMongo(&serverArgs.db.Mongo),
			log,
		)
	default:
		panic(base.Unreachable)
	}
}

func init() {
	defaultTimeout := 12 * time.Second
	base.AddServerFlags(serverCmd, &serverArgs.ServerOptions)
	base.AddServiceFlags(serverCmd, "seer", &serverArgs.seer, defaultTimeout)
	base.AddServiceFlags(serverCmd, "sources", &serverArgs.sources, defaultTimeout)
	base.AddServiceFlags(serverCmd, "slide-show", &serverArgs.slideshow, defaultTimeout)
	base.AddDatabaseFlags(serverCmd, &serverArgs.db)
	base.AddRedisFlags(serverCmd, &serverArgs.redis)
	serverCmd.PersistentFlags().IntVar(&serverArgs.concurrency, "concurrency", 1, "number of concurrent compile jobs")
	serverCmd.PersistentFlags().DurationVar(&serverArgs.saveInterval, "save-interval", 10*time.Minute, "save interval to produce an incomplete dataset that may be immediately used")
	serverCmd.PersistentFlags().BoolVarP(&serverArgs.compileOnStart, "compile-on-start", "c", false, "begin compilation on start")
	ConfigureCommand(serverCmd)
}
