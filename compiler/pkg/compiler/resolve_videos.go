package compiler

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/compiler/pkg/api"
	"github.com/thavlik/t4vd/compiler/pkg/datastore"
	seer "github.com/thavlik/t4vd/seer/pkg/api"
	"go.uber.org/zap"
)

// ResolveVideos resolved the Details field for each video in videos.
// If resolvedVideo is not nil, then each video's details will be sent to
// resolvedVideo as they are resolved.
func ResolveVideos(
	ctx context.Context,
	seerClient seer.Seer,
	ds datastore.DataStore,
	videos []*api.Video,
	resolvedVideo chan<- *api.Video,
	log *zap.Logger,
) error {
	if resolvedVideo != nil {
		defer close(resolvedVideo)
	}
	for _, video := range videos {
		resp, err := seerClient.GetVideoDetails(
			ctx,
			seer.GetVideoDetailsRequest{
				Input: video.ID,
			})
		if err != nil {
			return errors.Wrap(err, "seer.QueryVideoDetails")
		}
		video.Details = (*api.VideoDetails)(&resp.Details)
		if resolvedVideo != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case resolvedVideo <- video:
			}
		}
	}
	return nil
}
