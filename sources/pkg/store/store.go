package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/thavlik/t4vd/sources/pkg/api"
)

var ErrProjectNotFound = errors.New("project does not exist")

type Store interface {
	CreateProject(project *api.Project) error
	DeleteProject(id string) error
	ListProjects(ctx context.Context) ([]*api.Project, error)
	ListProjectsCreatedBy(ctx context.Context, userID string) ([]*api.Project, error)
	GetProject(ctx context.Context, id string) (*api.Project, error)
	GetProjectByName(ctx context.Context, name string) (*api.Project, error)
	GetProjectIDForGroup(ctx context.Context, groupID string) (projectID string, err error)
	GetProjectIDsForChannel(ctx context.Context, channelID string) ([]string, error)
	GetProjectIDsForPlaylist(ctx context.Context, playlistID string) ([]string, error)
	GetProjectIDsForVideo(ctx context.Context, videoID string) ([]string, error)
	AddChannel(projectID string, channel *api.Channel, blacklist bool, submitterID string) error
	AddPlaylist(projectID string, playlist *api.Playlist, blacklist bool, submitterID string) error
	AddVideo(projectID string, video *api.Video, blacklist bool, submitterID string) error
	ListChannels(ctx context.Context, projectID string) ([]*api.Channel, error)
	ListPlaylists(ctx context.Context, projectID string) ([]*api.Playlist, error)
	ListVideos(ctx context.Context, projectID string) ([]*api.Video, error)
	ListChannelIDs(ctx context.Context, projectID string, blacklist bool) ([]string, error)
	ListPlaylistIDs(ctx context.Context, projectID string, blacklist bool) ([]string, error)
	ListVideoIDs(ctx context.Context, projectID string, blacklist bool) ([]string, error)
	RemoveChannel(projectID string, channelID string, blacklist bool) error
	RemovePlaylist(projectID string, playlistID string, blacklist bool) error
	RemoveVideo(projectID string, videoID string, blacklist bool) error
}

func ScopedResourceID(projectID, id string) string {
	return fmt.Sprintf("%s:%s", projectID, id)
}
