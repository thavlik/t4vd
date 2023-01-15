package compiler

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/compiler/pkg/api"
	"github.com/thavlik/t4vd/compiler/pkg/datacache"
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
	dc datacache.DataCache,
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
	videoSources := make(map[string]*api.VideoSource)
	// Add input channels
	inputChannels, err := sourcesClient.ListChannels(
		ctx,
		sources.ListChannelsRequest{
			ProjectID: projectID,
		})
	if err != nil {
		return nil, errors.Wrap(err, "get input channels")
	}
	base.Progress(ctx, onProgress)
	log.Debug("queried channels",
		zap.Int("numChannels", len(inputChannels.Channels)))
	var blacklistChannelIDs []string
	for _, inputChannel := range inputChannels.Channels {
		if inputChannel.Blacklist {
			blacklistChannelIDs = append(blacklistChannelIDs, inputChannel.ID)
			continue
		}
		channelVideoIDs, err := GetChannelVideoIDs(ctx, seer, inputChannel.ID)
		if err != nil {
			return nil, errors.Wrap(err, "GetChannelVideoIDs")
		}
		log.Debug("queried whitelist channel videos",
			zap.String("channelID", inputChannel.ID),
			zap.Int("numVideos", len(channelVideoIDs)))
		for _, videoID := range channelVideoIDs {
			if _, ok := videoSources[videoID]; !ok {
				videoSources[videoID] = &api.VideoSource{
					Type:        "channel",
					ID:          inputChannel.ID,
					SubmitterID: inputChannel.SubmitterID,
					Submitted:   inputChannel.Submitted,
				}
			}
		}
		whitelist(videoIDs, channelVideoIDs)
	}
	base.Progress(ctx, onProgress)
	// Add input playlists
	inputPlaylists, err := sourcesClient.ListPlaylists(
		ctx,
		sources.ListPlaylistsRequest{
			ProjectID: projectID,
		})
	if err != nil {
		return nil, errors.Wrap(err, "get input playlists")
	}
	base.Progress(ctx, onProgress)
	log.Debug("queried playlists",
		zap.Int("numPlaylists", len(inputPlaylists.Playlists)))
	var blacklistPlaylistIDs []string
	for _, inputPlaylist := range inputPlaylists.Playlists {
		if inputPlaylist.Blacklist {
			blacklistPlaylistIDs = append(blacklistPlaylistIDs, inputPlaylist.ID)
			continue
		}
		playlistVideoIDs, err := GetPlaylistVideoIDs(ctx, seer, inputPlaylist.ID)
		if err != nil {
			return nil, errors.Wrap(err, "GetPlaylistVideoIDs")
		}
		log.Debug("queried whitelist playlist videos",
			zap.String("playlistID", inputPlaylist.ID),
			zap.Int("numVideos", len(playlistVideoIDs)))
		for _, videoID := range playlistVideoIDs {
			if _, ok := videoSources[videoID]; !ok {
				videoSources[videoID] = &api.VideoSource{
					Type:        "playlist",
					ID:          inputPlaylist.ID,
					SubmitterID: inputPlaylist.SubmitterID,
					Submitted:   inputPlaylist.Submitted,
				}
			}
		}
		whitelist(videoIDs, playlistVideoIDs)
	}
	base.Progress(ctx, onProgress)
	// Remove blacklist channels
	for _, blacklistChannelID := range blacklistChannelIDs {
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
	for _, blacklistPlaylistID := range blacklistPlaylistIDs {
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
	// Add input videos last so they are included
	// no matter the other blacklists.
	inputVideos, err := sourcesClient.ListVideos(
		ctx,
		sources.ListVideosRequest{
			ProjectID: projectID,
		})
	if err != nil {
		return nil, errors.Wrap(err, "get input videos")
	}
	var whitelistVideoIDs, blacklistVideoIDs []string
	for _, video := range inputVideos.Videos {
		if video.Blacklist {
			blacklistVideoIDs = append(blacklistVideoIDs, video.ID)
			continue
		}
		whitelistVideoIDs = append(whitelistVideoIDs, video.ID)
		videoSources[video.ID] = &api.VideoSource{
			SubmitterID: video.SubmitterID,
			Submitted:   video.Submitted,
		}
	}
	blacklist(videoIDs, blacklistVideoIDs)
	whitelist(videoIDs, whitelistVideoIDs)
	base.Progress(ctx, onProgress)
	compiled := make([]*api.Video, len(videoIDs))
	i := 0
	for videoID := range videoIDs {
		source, ok := videoSources[videoID]
		if !ok {
			return nil, errors.New("failed sanity check: video source not found")
		}
		compiled[i] = &api.Video{
			ID:     videoID,
			Source: source,
		}
		i++
	}

	// we can now cache the video IDs for quick lookup
	if err := dc.Set(projectID, flatten(videoIDs)); err != nil {
		return nil, errors.Wrap(err, "datacache.Set")
	}

	log.Debug("resolving videos",
		base.Elapsed(start),
		zap.Int("count", len(compiled)))
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
					if err := dc.Add(projectID, video.ID); err != nil {
						log.Warn("failed to add video to datacache", zap.Error(err))
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
	if err := ResolveVideos(
		ctx,
		seer,
		ds,
		compiled,
		resolvedVideo,
		log,
	); err != nil {
		return nil, errors.Wrap(err, "ResolveVideos")
	}
	base.Progress(ctx, onProgress)
	log.Debug("resolved videos",
		base.Elapsed(start),
		zap.Int("count", len(compiled)))
	// save completed dataset
	if len(compiled) == 0 {
		return nil, ErrNoVideos
	}
	log.Debug("saving complete dataset",
		base.Elapsed(start),
		zap.Int("numVideos", len(compiled)))
	dataset, err := ds.SaveDataset(
		context.Background(), // do not allow interruption
		projectID,
		compiled,
		true,
		start,
	)
	if err != nil {
		return nil, errors.Wrap(err, "datastore.SaveDataset")
	}
	log.Debug("finished compiling dataset",
		base.Elapsed(start),
		zap.String("projectID", projectID),
		zap.Int("numVideos", len(compiled)),
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
