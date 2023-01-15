package server

import (
	"context"
	"math"
	"math/rand"
	"time"

	"github.com/mroth/weightedrand/v2"
	"github.com/thavlik/t4vd/base/pkg/base"
	compiler "github.com/thavlik/t4vd/compiler/pkg/api"
	"go.uber.org/zap"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/slideshow/pkg/api"
	"github.com/thavlik/t4vd/slideshow/pkg/imgcache"
)

var frameTime time.Duration = 33333333

func genRandomMarker(
	ctx context.Context,
	imgCache imgcache.ImgCache,
	compilerClient compiler.Compiler,
	projectID string,
	bucket string,
	log *zap.Logger,
) (*api.Marker, error) {
	log = log.With(zap.String("projectID", projectID))
	log.Debug("generating random marker")
	start := time.Now()
	ds, err := compilerClient.GetDataset(ctx, compiler.GetDatasetRequest{
		ProjectID: projectID,
	})
	if err != nil {
		log.Warn("failed to generate marker", zap.Error(err))
		return nil, err
	}
	log.Debug("retrieved dataset",
		base.Elapsed(start),
		zap.String("id", ds.ID),
		zap.Int("numVideos", len(ds.Videos)))
	videoIDs, err := listCachedVideoIDs(bucket)
	if err != nil {
		return nil, errors.Wrap(err, "listCachedVideoIDs")
	}
	videos := keepOnly(ds.Videos, videoIDs)
	log.Debug("calculated videos",
		base.Elapsed(start),
		zap.Int("numVideos", len(videos)))
	if len(videos) == 0 {
		return nil, errors.New("no videos cached")
	}
	choices := make([]weightedrand.Choice[*compiler.Video, int64], len(videos))
	for i, video := range videos {
		choices[i] = weightedrand.NewChoice(video, video.Details.Duration)
	}
	chooser, err := weightedrand.NewChooser(choices...)
	if err != nil {
		return nil, errors.Wrap(err, "weightedrand.NewChooser")
	}
	video := chooser.Pick()
	t := time.Duration(math.Floor(rand.Float64() * float64(video.Details.Duration*int64(time.Second))))
	// cache the image data
	//t -= t % frameTime // round to nearest frame
	if _, err := getFrame(
		ctx,
		imgCache,
		bucket,
		video.ID,
		t,
		log,
	); err != nil {
		return nil, errors.Wrap(err, "getFrame")
	}
	log.Debug("generated marker", base.Elapsed(start))
	return &api.Marker{
		VideoID: video.ID,
		Time:    int64(t),
	}, nil
}
