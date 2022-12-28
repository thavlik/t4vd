package infocache

import (
	"context"
	"errors"
	"time"

	"github.com/thavlik/bjjvb/seer/pkg/api"
)

var (
	ErrCacheUnavailable = errors.New("cache unavailable")
	ErrCacheExpired     = errors.New("cache expired")
	CacheRecency        = 96 * time.Hour
)

type InfoCache interface {
	GetVideo(ctx context.Context, videoID string) (*api.VideoDetails, error)
	SetVideo(video *api.VideoDetails) error

	GetChannel(ctx context.Context, channelID string) (*api.ChannelDetails, error)
	SetChannel(*api.ChannelDetails) error
	GetChannelVideoIDs(ctx context.Context, channelID string) (videoIDs []string, timestamp time.Time, err error)
	SetChannelVideoIDs(channelID string, videoIDs []string, timestamp time.Time) error

	GetPlaylist(ctx context.Context, playlistID string) (*api.PlaylistDetails, error)
	SetPlaylist(*api.PlaylistDetails) error
	GetPlaylistVideoIDs(ctx context.Context, channelID string) (videoIDs []string, timestamp time.Time, err error)
	SetPlaylistVideoIDs(channelID string, videoIDs []string, timestamp time.Time) error
}
