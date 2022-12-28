package datastore

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/compiler/pkg/api"
	seer "github.com/thavlik/bjjvb/seer/pkg/api"
	"go.uber.org/zap"
)

var errSanityCheck = errors.New("sanity check failed")

// ResolveVideos transforms a list of videoIDs into full
// json using the seer microservice.
func ResolveVideos(
	ctx context.Context,
	seerClient seer.Seer,
	ds DataStore,
	videoIDs []string,
	resolvedVideo chan<- *api.Video,
	log *zap.Logger,
) ([]*api.Video, error) {
	if resolvedVideo != nil {
		defer close(resolvedVideo)
	}
	videos := make([]*api.Video, len(videoIDs))
	for i, videoID := range videoIDs {
		if videoID == "" {
			return nil, errors.Wrap(errSanityCheck, "videoID is empty")
		}
		video, err := ds.GetCachedVideo(ctx, videoID)
		if err == ErrVideoNotCached {
			// Get video info from seer
			resp, err := seerClient.GetVideoDetails(ctx, seer.GetVideoDetailsRequest{
				Input: videoID,
			})
			if err != nil {
				return nil, errors.Wrap(err, "seer.QueryVideoDetails")
			}
			video = &api.Video{
				ID:          resp.Details.ID,
				Title:       resp.Details.Title,
				Description: resp.Details.Description,
				Thumbnail:   resp.Details.Thumbnail,
				UploadDate:  resp.Details.UploadDate,
				Uploader:    resp.Details.Uploader,
				UploaderID:  resp.Details.UploaderID,
				Channel:     resp.Details.Channel,
				ChannelID:   resp.Details.ChannelID,
				Duration:    resp.Details.Duration,
				ViewCount:   resp.Details.ViewCount,
				Width:       resp.Details.Width,
				Height:      resp.Details.Height,
				FPS:         resp.Details.FPS,
			}
			if resolvedVideo != nil {
				select {
				case <-ctx.Done():
					return nil, errors.Wrap(ctx.Err(), "context")
				case resolvedVideo <- video:
				}
			}
			if err := ds.CacheVideo(ctx, video); err != nil {
				return nil, errors.Wrap(err, "datastore.CacheVideo")
			}
		} else if err != nil {
			return nil, errors.Wrap(err, "GetCachedVideo")
		} else if resolvedVideo != nil {
			select {
			case <-ctx.Done():
				return nil, errors.Wrap(ctx.Err(), "context")
			case resolvedVideo <- video:
			}
		}
		if video.ID == "" {
			return nil, errors.Wrap(errSanityCheck, "sanity check failed: video has no id")
		}
		videos[i] = video
	}
	return videos, nil
}
