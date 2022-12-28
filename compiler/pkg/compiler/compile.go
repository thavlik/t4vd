package compiler

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/compiler/pkg/api"
	"github.com/thavlik/t4vd/compiler/pkg/datastore"
	seer "github.com/thavlik/t4vd/seer/pkg/api"
	sources "github.com/thavlik/t4vd/sources/pkg/api"
	"go.uber.org/zap"
)

var ErrNoVideos = errors.New("dataset has no videos")

func Compile(
	ctx context.Context,
	projectID string,
	sourcesClient sources.Sources,
	seer seer.Seer,
	ds datastore.DataStore,
	saveInterval time.Duration,
	saved chan<- *api.Dataset,
	onProgress chan<- struct{},
	log *zap.Logger,
) (*api.Dataset, error) {
	base.Progress(ctx, onProgress)
	defer base.Progress(ctx, onProgress)
	start := time.Now()
	log = log.With(zap.String("projectID", projectID))
	log.Debug("compiling dataset")
	videoIDs := make(map[string]struct{})
	// Add input channels
	inputChannels, err := sourcesClient.ListChannelIDs(
		ctx,
		sources.ListChannelIDsRequest{
			ProjectID: projectID,
			Blacklist: false,
		})
	if err != nil {
		return nil, errors.Wrap(err, "get input channels")
	}
	base.Progress(ctx, onProgress)
	log.Debug("queried whitelist channels",
		zap.Int("numChannels", len(inputChannels.IDs)))
	for _, inputChannelID := range inputChannels.IDs {
		channelVideoIDs, err := GetChannelVideoIDs(ctx, seer, inputChannelID)
		if err != nil {
			return nil, errors.Wrap(err, "GetChannelVideoIDs")
		}
		log.Debug("queried whitelist channel videos",
			zap.String("channelID", inputChannelID),
			zap.Int("numVideos", len(channelVideoIDs)))
		whitelist(videoIDs, channelVideoIDs)
	}
	base.Progress(ctx, onProgress)
	// Add input playlists
	inputPlaylists, err := sourcesClient.ListPlaylistIDs(
		ctx,
		sources.ListPlaylistIDsRequest{
			ProjectID: projectID,
			Blacklist: false,
		})
	if err != nil {
		return nil, errors.Wrap(err, "get input playlists")
	}
	base.Progress(ctx, onProgress)
	log.Debug("queried whitelist playlists",
		zap.Int("numPlaylists", len(inputPlaylists.IDs)))
	for _, inputPlaylistID := range inputPlaylists.IDs {
		playlistVideoIDs, err := GetPlaylistVideoIDs(ctx, seer, inputPlaylistID)
		if err != nil {
			return nil, errors.Wrap(err, "GetPlaylistVideoIDs")
		}
		log.Debug("queried whitelist playlist videos",
			zap.String("playlistID", inputPlaylistID),
			zap.Int("numVideos", len(playlistVideoIDs)))
		whitelist(videoIDs, playlistVideoIDs)
	}
	base.Progress(ctx, onProgress)
	// Remove blacklist channels
	blacklistChannels, err := sourcesClient.ListChannelIDs(
		ctx,
		sources.ListChannelIDsRequest{
			ProjectID: projectID,
			Blacklist: true,
		})
	if err != nil {
		return nil, errors.Wrap(err, "get blacklist channels")
	}
	base.Progress(ctx, onProgress)
	log.Debug("queried blacklist channels",
		zap.Int("numChannels", len(blacklistChannels.IDs)))
	for _, blacklistChannelID := range blacklistChannels.IDs {
		channelVideoIDs, err := GetChannelVideoIDs(ctx, seer, blacklistChannelID)
		if err != nil {
			return nil, errors.Wrap(err, "GetChannelVideoIDs")
		}
		log.Debug("queried blacklist channel videos",
			zap.String("channelID", blacklistChannelID),
			zap.Int("numVideos", len(channelVideoIDs)))
		blacklist(videoIDs, channelVideoIDs)
	}
	base.Progress(ctx, onProgress)
	// Remove blacklist playlists
	blacklistPlaylists, err := sourcesClient.ListPlaylistIDs(
		ctx,
		sources.ListPlaylistIDsRequest{
			ProjectID: projectID,
			Blacklist: true,
		})
	if err != nil {
		return nil, errors.Wrap(err, "get blacklist playlists")
	}
	base.Progress(ctx, onProgress)
	log.Debug("queried blacklist playlists",
		zap.Int("numPlaylists", len(blacklistPlaylists.IDs)))
	for _, blacklistPlaylistID := range blacklistPlaylists.IDs {
		playlistVideoIDs, err := GetPlaylistVideoIDs(ctx, seer, blacklistPlaylistID)
		if err != nil {
			return nil, errors.Wrap(err, "GetPlaylistVideoIDs")
		}
		log.Debug("queried blacklist playlist videos",
			zap.String("playlistID", blacklistPlaylistID),
			zap.Int("numVideos", len(playlistVideoIDs)))
		blacklist(videoIDs, playlistVideoIDs)
	}
	base.Progress(ctx, onProgress)
	// Remove blacklist videos
	blacklistVideos, err := sourcesClient.ListVideoIDs(
		ctx,
		sources.ListVideoIDsRequest{
			ProjectID: projectID,
			Blacklist: true,
		})
	if err != nil {
		return nil, errors.Wrap(err, "get blacklist videos")
	}
	base.Progress(ctx, onProgress)
	log.Debug("queried blacklist videos",
		zap.Int("numVideos", len(blacklistVideos.IDs)))
	blacklist(videoIDs, blacklistVideos.IDs)
	// Add input videos last so they are included
	// no matter the other blacklists.
	inputVideos, err := sourcesClient.ListVideoIDs(
		ctx,
		sources.ListVideoIDsRequest{
			ProjectID: projectID,
			Blacklist: false,
		})
	if err != nil {
		return nil, errors.Wrap(err, "get input videos")
	}
	base.Progress(ctx, onProgress)
	log.Debug("queried whitelist videos",
		zap.Int("numVideos", len(inputVideos.IDs)))
	whitelist(videoIDs, inputVideos.IDs)
	flattened := flatten(videoIDs)
	log.Debug("resolving videos",
		base.Elapsed(start),
		zap.Int("count", len(flattened)))
	var resolvedVideo chan *api.Video
	if saveInterval != 0 {
		resolvedVideo = make(chan *api.Video, 1)
		stop := make(chan struct{}, 1)
		stopped := make(chan struct{})
		defer func() {
			stop <- struct{}{}
			<-stopped
		}()
		go func() {
			defer func() {
				stopped <- struct{}{}
			}()
			var videos []*api.Video
			if saved != nil {
				defer close(saved)
			}
			var save <-chan time.Time
			resetTimer := func() { save = time.After(saveInterval) }
			resetTimer()
			for {
				select {
				case <-stop:
					return
				case video, ok := <-resolvedVideo:
					if !ok {
						return
					}
					videos = append(videos, video)
					base.Progress(ctx, onProgress)
					continue
				case _, ok := <-save:
					if !ok {
						return
					}
					if len(videos) == 0 {
						// no videos available yet
						log.Warn("no videos available yet")
						resetTimer()
						continue
					}
					base.Progress(ctx, onProgress)
					// autosaved datasets have complete=false
					dataset, err := ds.SaveDataset(
						ctx,
						projectID,
						videos,
						false,
						start,
					)
					if err != nil {
						log.Error("failed to autosave dataset",
							zap.Error(errors.Wrap(err, "datastore.SaveDataset")))
						resetTimer()
						continue
					}
					log.Debug("autosaved compiling dataset",
						base.Elapsed(start),
						zap.Int("numVideos", len(videos)),
						zap.String("dataset.ID", dataset.ID))
					if saved != nil {
						select {
						case <-stop:
							return
						case saved <- dataset:
						}
					}
					resetTimer()
				}
			}
		}()
	}
	base.Progress(ctx, onProgress)
	videos, err := datastore.ResolveVideos(
		ctx,
		seer,
		ds,
		flattened,
		resolvedVideo,
		log,
	)
	if err != nil {
		return nil, errors.Wrap(err, "ResolveVideos")
	}
	base.Progress(ctx, onProgress)
	log.Debug("resolved videos",
		base.Elapsed(start),
		zap.Int("count", len(videos)))
	// save completed dataset
	if len(videos) == 0 {
		return nil, ErrNoVideos
	}
	log.Debug("saving complete dataset",
		base.Elapsed(start),
		zap.Int("numVideos", len(videos)))
	dataset, err := ds.SaveDataset(context.Background(), projectID, videos, true, start)
	if err != nil {
		return nil, errors.Wrap(err, "datastore.SaveDataset")
	}
	log.Debug("finished compiling dataset",
		base.Elapsed(start),
		zap.String("projectID", projectID),
		zap.Int("numVideos", len(flattened)),
		zap.String("dataset.ID", dataset.ID))
	return dataset, nil
}

// GetChannelVideoIDs returns all video IDs for a channel
func GetChannelVideoIDs(
	ctx context.Context,
	seerClient seer.Seer,
	channelID string,
) ([]string, error) {
	resp, err := seerClient.GetChannelVideoIDs(
		ctx,
		seer.GetChannelVideoIDsRequest{
			ID: channelID,
		})
	if err != nil {
		return nil, errors.Wrap(err, "seer")
	}
	return resp.VideoIDs, nil
}

// GetPlaylistVideoIDs returns all video IDs for a playlist
func GetPlaylistVideoIDs(
	ctx context.Context,
	seerClient seer.Seer,
	playlistID string,
) ([]string, error) {
	resp, err := seerClient.GetPlaylistVideoIDs(
		ctx,
		seer.GetPlaylistVideoIDsRequest{
			ID: playlistID,
		})
	if err != nil {
		return nil, errors.Wrap(err, "seer")
	}
	return resp.VideoIDs, nil
}

func whitelist(m map[string]struct{}, ids []string) {
	for _, id := range ids {
		m[id] = struct{}{}
	}
}

func blacklist(input map[string]struct{}, blacklist []string) {
	for _, k := range blacklist {
		delete(input, k)
	}
}

func flatten(m map[string]struct{}) []string {
	var result []string
	for k := range m {
		result = append(result, k)
	}
	return result
}
