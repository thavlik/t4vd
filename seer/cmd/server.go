package main

import (
	"errors"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/spf13/cobra"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/pubsub"
	memory_pubsub "github.com/thavlik/t4vd/base/pkg/pubsub/memory"
	redis_pubsub "github.com/thavlik/t4vd/base/pkg/pubsub/redis"
	"github.com/thavlik/t4vd/base/pkg/scheduler"
	memory_scheduler "github.com/thavlik/t4vd/base/pkg/scheduler/memory"
	redis_scheduler "github.com/thavlik/t4vd/base/pkg/scheduler/redis"
	hound "github.com/thavlik/t4vd/hound/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/cachedset"
	redis_cachedset "github.com/thavlik/t4vd/seer/pkg/cachedset/redis"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	mongo_infocache "github.com/thavlik/t4vd/seer/pkg/infocache/mongo"
	postgres_infocache "github.com/thavlik/t4vd/seer/pkg/infocache/postgres"
	"github.com/thavlik/t4vd/seer/pkg/server"
	"github.com/thavlik/t4vd/seer/pkg/thumbcache"
	s3_thumbcache "github.com/thavlik/t4vd/seer/pkg/thumbcache/s3"
	"github.com/thavlik/t4vd/seer/pkg/vidcache"
	s3_vidcache "github.com/thavlik/t4vd/seer/pkg/vidcache/s3"
	"go.uber.org/zap"
)

var serverArgs struct {
	base.ServerOptions
	redis            base.RedisOptions
	db               base.DatabaseOptions
	hound            base.ServiceOptions
	videoBucket      string
	thumbnailBucket  string
	videoFormat      string
	includeAudio     bool
	concurrency      int
	disableDownloads bool
}

var serverCmd = &cobra.Command{
	Use: "server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.ServerEnv(&serverArgs.ServerOptions)
		base.DatabaseEnv(&serverArgs.db, true)
		base.RedisEnv(&serverArgs.redis, false)
		base.ServiceEnv("hound", &serverArgs.hound)
		base.CheckEnv("VIDEO_BUCKET", &serverArgs.videoBucket)
		base.CheckEnv("VIDEO_FORMAT", &serverArgs.videoFormat)
		base.CheckEnvBool("INCLUDE_AUDIO", &serverArgs.includeAudio)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log := base.DefaultLog
		go base.RunMetrics(serverArgs.MetricsPort, log)
		redis := initRedis()
		return server.Entry(
			serverArgs.Port,
			initPubSub(redis, log),
			initScheduler(redis, "dlsched"),
			initScheduler(redis, "qrysched"),
			initInfoCache(&serverArgs.db),
			initVidCache(log),
			initThumbCache(log),
			initCachedSet(redis, log),
			hound.NewHoundClientFromOptions(serverArgs.hound),
			serverArgs.videoFormat,
			serverArgs.includeAudio,
			serverArgs.concurrency,
			serverArgs.disableDownloads,
			log,
		)
	},
}

func initRedis() *redis.Client {
	if serverArgs.redis.IsSet() {
		return base.ConnectRedis(&serverArgs.redis)
	}
	return nil
}

func initCachedSet(redis *redis.Client, log *zap.Logger) cachedset.CachedSet {
	if redis != nil {
		return redis_cachedset.NewRedisCachedSet(redis, log)
	}
	panic(errors.New("missing cached set source"))
}

func initPubSub(
	redis *redis.Client,
	log *zap.Logger,
) pubsub.PubSub {
	if redis != nil {
		return redis_pubsub.NewRedisPubSub(
			redis,
			log,
		)
	}
	return memory_pubsub.NewMemoryPubSub(log)
}

func initScheduler(
	redis *redis.Client,
	name string,
) scheduler.Scheduler {
	if redis != nil {
		return redis_scheduler.NewRedisScheduler(
			redis,
			name,
			25*time.Second,
		)
	}
	return memory_scheduler.NewMemoryScheduler()
}

func initVidCache(log *zap.Logger) vidcache.VidCache {
	if serverArgs.videoBucket != "" {
		return s3_vidcache.NewS3VidCache(serverArgs.videoBucket, log)
	} else {
		panic(errors.New("missing video cache source"))
	}
}

func initThumbCache(log *zap.Logger) thumbcache.ThumbCache {
	if serverArgs.thumbnailBucket != "" {
		return s3_thumbcache.NewS3ThumbCache(serverArgs.thumbnailBucket, log)
	} else {
		panic(errors.New("missing thumbnail cache source"))
	}
}

func initInfoCache(opts *base.DatabaseOptions) infocache.InfoCache {
	switch opts.Driver {
	case "":
		panic("missing --db-driver")
	case base.MongoDriver:
		return mongo_infocache.NewMongoInfoCache(
			base.ConnectMongo(&opts.Mongo))
	case base.PostgresDriver:
		return postgres_infocache.NewPostgresInfoCache(
			base.ConnectPostgres(&opts.Postgres))
	default:
		panic(base.Unreachable)
	}
}

func init() {
	base.AddRedisFlags(serverCmd, &serverArgs.redis)
	base.AddServerFlags(serverCmd, &serverArgs.ServerOptions)
	base.AddDatabaseFlags(serverCmd, &serverArgs.db)
	base.AddServiceFlags(serverCmd, "hound", &serverArgs.hound, 10*time.Second)
	serverCmd.PersistentFlags().IntVar(&serverArgs.concurrency, "concurrency", 1, "number of concurrent youtube queries (best set to 1 and increase # replicas)")
	serverCmd.PersistentFlags().StringVar(&serverArgs.videoBucket, "video-bucket", "", "full length video cache bucket name")
	serverCmd.PersistentFlags().StringVar(&serverArgs.thumbnailBucket, "thumbnail-bucket", "", "thumbnail cache bucket name")
	serverCmd.PersistentFlags().StringVar(&serverArgs.videoFormat, "video-format", "webm", "download video format")
	serverCmd.PersistentFlags().BoolVar(&serverArgs.includeAudio, "include-audio", false, "download audio")
	serverCmd.PersistentFlags().BoolVar(&serverArgs.disableDownloads, "disable-downloads", false, "disable all downloads from youtube (info queries still allowed)")
	ConfigureCommand(serverCmd)
}
