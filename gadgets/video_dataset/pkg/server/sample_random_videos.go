package server

import (
	"context"

	"github.com/mroth/weightedrand/v2"
	"github.com/pkg/errors"
	compiler "github.com/thavlik/t4vd/compiler/pkg/api"
)

func sampleRandomVideos(
	ctx context.Context,
	compilerClient compiler.Compiler,
	projectID string,
	batchSize int,
) ([]*compiler.Video, error) {
	ds, err := compilerClient.GetDataset(
		ctx,
		compiler.GetDatasetRequest{
			ProjectID: projectID,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "compiler.GetDataset")
	}
	choices := make([]weightedrand.Choice[*compiler.Video, int64], len(ds.Videos))
	for i, video := range ds.Videos {
		choices[i] = weightedrand.NewChoice(video, video.Details.Duration)
	}
	chooser, err := weightedrand.NewChooser(choices...)
	if err != nil {
		return nil, errors.Wrap(err, "weightedrand.NewChooser")
	}
	videos := make([]*compiler.Video, batchSize)
	for i := 0; i < batchSize; i++ {
		videos[i] = chooser.Pick()
	}
	return videos, nil
}
