package infocache

import (
	"context"
	"errors"
	"time"

	"github.com/thavlik/t4vd/seer/pkg/api"
)

var (
	ErrCacheUnavailable = errors.New("cache unavailable")
	ErrCacheExpired     = errors.New("cache expired")
	CacheRecency        = 96 * time.Hour
)

type InfoCache interface {
	GetBulkVideos(ctx context.Context, videoIDs []string) ([]*api.VideoDetails, error)
	GetBulkPlaylists(ctx context.Context, playlistIDs []string) ([]*api.PlaylistDetails, error)
	GetBulkChannels(ctx context.Context, channelIDs []string) ([]*api.ChannelDetails, error)

	GetVideo(ctx context.Context, videoID string) (*api.VideoDetails, error)
	SetVideo(video *api.VideoDetails) error
	IsVideoRecent(videoID string) (bool, error)

	GetChannel(ctx context.Context, channelID string) (*api.ChannelDetails, error)
	SetChannel(*api.ChannelDetails) error
	GetChannelVideoIDs(ctx context.Context, channelID string) (videoIDs []string, timestamp time.Time, err error)
	SetChannelVideoIDs(channelID string, videoIDs []string, timestamp time.Time) error
	IsChannelRecent(channelID string) (bool, error)

	GetPlaylist(ctx context.Context, playlistID string) (*api.PlaylistDetails, error)
	SetPlaylist(*api.PlaylistDetails) error
	GetPlaylistVideoIDs(ctx context.Context, channelID string) (videoIDs []string, timestamp time.Time, err error)
	SetPlaylistVideoIDs(channelID string, videoIDs []string, timestamp time.Time) error
	IsPlaylistRecent(playlistID string) (bool, error)
}
