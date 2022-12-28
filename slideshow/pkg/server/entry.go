package server

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/base/pkg/base"
	compiler "github.com/thavlik/bjjvb/compiler/pkg/api"
	"github.com/thavlik/bjjvb/slideshow/pkg/api"
	"github.com/thavlik/bjjvb/slideshow/pkg/imgcache"
	redis_markercache "github.com/thavlik/bjjvb/slideshow/pkg/markercache/redis"

	"github.com/thavlik/bjjvb/slideshow/pkg/markercache"
	sources "github.com/thavlik/bjjvb/sources/pkg/api"
	"go.uber.org/zap"
)

var ErrMissingProjectID = errors.New("missing project id")

func Entry(
	port int,
	bucket string,
	imgCache imgcache.ImgCache,
	redis *redis.Client,
	compiler compiler.Compiler,
	sourcesClient sources.Sources,
	log *zap.Logger,
) (err error) {
	markerCache := markerCacheClient(
		redis,
		func(projectID string) (*api.Marker, error) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			return genRandomMarker(
				ctx,
				imgCache,
				compiler,
				projectID,
				bucket,
				log,
			)
		},
		log,
	)
	markerCache.Start()
	s := NewServer(
		bucket,
		imgCache,
		markerCache,
		compiler,
		log,
	)
	defer s.markerCache.Close()
	resp, err := sourcesClient.ListProjects(
		context.Background(),
		sources.ListProjectsRequest{})
	for _, project := range resp.Projects {
		if err := s.markerCache.Queue(project.ID); err != nil {
			panic(err)
		}
	}
	base.SignalReady(log)
	return s.ListenAndServe(port)
}

func markerCacheClient(
	redisClient *redis.Client,
	genMarker markercache.GenMarkerFunc,
	log *zap.Logger,
) markercache.MarkerCache {
	if redisClient == nil {
		panic("redis is currently required")
	}
	cacheSize := base.CheckEnvInt64("MARKER_CACHE_SIZE", nil)
	if cacheSize == 0 {
		cacheSize = 16
	}
	return redis_markercache.NewRedisMarkerCache(
		redisClient,
		genMarker,
		cacheSize,
		log,
	)
}
