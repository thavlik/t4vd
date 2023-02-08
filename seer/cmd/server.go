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
	"github.com/thavlik/t4vd/seer/pkg/ytdl"
	"go.uber.org/zap"
)

var serverArgs struct {
	base.ServerOptions
	redis             base.RedisOptions
	db                base.DatabaseOptions
	hound             base.ServiceOptions
	videoBucket       string
	thumbnailBucket   string
	videoFormat       string
	audioFormat       string
	audioChannelCount int
	audioSampleRate   int
	skipAudio         bool
	skipVideo         bool
	concurrency       int
	disableDownloads  bool
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
		base.CheckEnvInt("CONCURRENCY", &serverArgs.concurrency)
		base.CheckEnvBool("DISABLE_DOWNLOADS", &serverArgs.disableDownloads)
		base.CheckEnv("THUMBNAIL_BUCKET", &serverArgs.thumbnailBucket)
		base.CheckEnvBool("SKIP_AUDIO", &serverArgs.skipAudio)
		base.CheckEnvBool("SKIP_VIDEO", &serverArgs.skipVideo)
		base.CheckEnv("AUDIO_FORMAT", &serverArgs.audioFormat)
		base.CheckEnvInt("AUDIO_CHANNEL_COUNT", &serverArgs.audioChannelCount)
		base.CheckEnvInt("AUDIO_SAMPLE_RATE", &serverArgs.audioSampleRate)
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
			&ytdl.Options{
				VideoFormat:       serverArgs.videoFormat,
				AudioFormat:       serverArgs.audioFormat,
				SkipAudio:         serverArgs.skipAudio,
				SkipVideo:         serverArgs.skipVideo,
				AudioChannelCount: serverArgs.audioChannelCount,
				AudioSampleRate:   serverArgs.audioSampleRate,
			},
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
		return s3_vidcache.NewS3VidCache(
			serverArgs.videoBucket,
			serverArgs.videoFormat,
			log,
		)
	} else {
		panic(errors.New("missing video cache source"))
	}
}

func initThumbCache(log *zap.Logger) thumbcache.ThumbCache {
	if serverArgs.thumbnailBucket != "" {
		return s3_thumbcache.NewS3ThumbCache(
			serverArgs.thumbnailBucket,
			log,
		)
	} else {
		// don't cache thumbnails
		return nil
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
	serverCmd.PersistentFlags().StringVar(&serverArgs.thumbnailBucket, "thumbnail-bucket", "", "thumbnail cache bucket name (optional)")
	serverCmd.PersistentFlags().StringVar(&serverArgs.videoFormat, "video-format", "webm", "download video format")
	serverCmd.PersistentFlags().StringVar(&serverArgs.audioFormat, "audio-format", "webm", "download audio format")
	serverCmd.PersistentFlags().IntVar(&serverArgs.audioChannelCount, "audio-channel-count", 1, "download audio channel count (1 = mono, 2 = stereo)")
	serverCmd.PersistentFlags().IntVar(&serverArgs.audioSampleRate, "audio-sample-rate", 44100, "download audio sample rate")
	serverCmd.PersistentFlags().BoolVar(&serverArgs.skipAudio, "skip-audio", false, "skip downloading audio")
	serverCmd.PersistentFlags().BoolVar(&serverArgs.skipVideo, "skip-video", false, "skip downloading video")
	serverCmd.PersistentFlags().BoolVar(&serverArgs.disableDownloads, "disable-downloads", false, "disable all downloads from youtube (info queries still allowed)")
	ConfigureCommand(serverCmd)
}
