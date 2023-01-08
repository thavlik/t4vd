package definitions

type Sources interface {
	CreateProject(Project) Project
	DeleteProject(DeleteProject) Void
	ListProjects(ListProjectsRequest) ListProjectsResponse
	GetProject(GetProject) Project
	GetProjectByName(GetProjectByName) Project

	GetProjectIDsForChannel(GetProjectIDsForChannelRequest) GetProjectIDsForChannelResponse
	GetProjectIDsForPlaylist(GetProjectIDsForPlaylistRequest) GetProjectIDsForPlaylistResponse
	GetProjectIDsForVideo(GetProjectIDsForVideoRequest) GetProjectIDsForVideoResponse

	AddChannel(AddChannelRequest) Channel
	AddPlaylist(AddPlaylistRequest) Playlist
	AddVideo(AddVideoRequest) Video

	ListChannels(ListChannelsRequest) ListChannelsResponse
	ListPlaylists(ListPlaylistsRequest) ListPlaylistsResponse
	ListVideos(ListVideosRequest) ListVideosResponse

	ListChannelIDs(ListChannelIDsRequest) ListChannelIDsResponse
	ListPlaylistIDs(ListPlaylistIDsRequest) ListPlaylistIDsResponse
	ListVideoIDs(ListVideoIDsRequest) ListVideoIDsResponse

	RemoveChannel(RemoveChannelRequest) Void
	RemovePlaylist(RemovePlaylistRequest) Void
	RemoveVideo(RemoveVideoRequest) Void
}

type Void struct{}

type GetProjectIDsForChannelRequest struct {
	ChannelID string `json:"channelID"`
}

type GetProjectIDsForChannelResponse struct {
	ProjectIDs []string `json:"projectIDs"`
}

type GetProjectIDsForPlaylistRequest struct {
	PlaylistID string `json:"playlistID"`
}

type GetProjectIDsForPlaylistResponse struct {
	ProjectIDs []string `json:"projectIDs"`
}

type GetProjectIDsForVideoRequest struct {
	VideoID string `json:"videoID"`
}

type GetProjectIDsForVideoResponse struct {
	ProjectIDs []string `json:"projectIDs"`
}

type GetProject struct {
	ID string `json:"id"`
}

type GetProjectByName struct {
	Name string `json:"name"`
}

type ListProjectsRequest struct {
	CreatedByUserID string `json:"createdByUserID"`
	VisibleToUserID string `json:"visibleToUserID"`
}

type ListProjectsResponse struct {
	Projects []*Project `json:"projects"`
}

type Collaborator struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatorID string `json:"creatorID"`
	GroupID   string `json:"groupID"`
}

type DeleteProject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AddChannelRequest struct {
	ProjectID   string `json:"projectID"`
	Input       string `json:"input"`
	Blacklist   bool   `json:"blacklist"`
	SubmitterID string `json:"submitterID"`
}

type RemoveChannelRequest struct {
	ProjectID string `json:"projectID"`
	ID        string `json:"ID"`
	Blacklist bool   `json:"blacklist"`
}

type AddPlaylistRequest struct {
	ProjectID   string `json:"projectID"`
	Input       string `json:"input"`
	Blacklist   bool   `json:"blacklist"`
	SubmitterID string `json:"submitterID"`
}

type RemovePlaylistRequest struct {
	ProjectID string `json:"projectID"`
	ID        string `json:"ID"`
	Blacklist bool   `json:"blacklist"`
}

type AddVideoRequest struct {
	ProjectID   string `json:"projectID"`
	Input       string `json:"input"`
	Blacklist   bool   `json:"blacklist"`
	SubmitterID string `json:"submitterID"`
}

type RemoveVideoRequest struct {
	ProjectID string `json:"projectID"`
	ID        string `json:"ID"`
	Blacklist bool   `json:"blacklist"`
}

type ListChannelsRequest struct {
	ProjectID string `json:"projectID"`
	Blacklist bool   `json:"blacklist"`
}

type Channel struct {
	ID        string `json:"id"`
	Blacklist bool   `json:"blacklist"`
}

type ListChannelsResponse struct {
	Channels []*Channel `json:"channels"`
}

type ListPlaylistsRequest struct {
	ProjectID string `json:"projectID"`
	Blacklist bool   `json:"blacklist"`
}

type Playlist struct {
	ID        string `json:"id"`
	Blacklist bool   `json:"blacklist"`
}

type ListPlaylistsResponse struct {
	Playlists []*Playlist `json:"playlists"`
}

type ListVideosRequest struct {
	ProjectID string `json:"projectID"`
	Blacklist bool   `json:"blacklist"`
}

type Video struct {
	ID        string `json:"id"`
	Blacklist bool   `json:"blacklist"`
}

type ListVideosResponse struct {
	Videos []*Video `json:"videos"`
}

type ListChannelIDsRequest struct {
	ProjectID string `json:"projectID"`
	Blacklist bool   `json:"blacklist"`
}

type ListChannelIDsResponse struct {
	IDs []string `json:"IDs"`
}

type ListPlaylistIDsRequest struct {
	ProjectID string `json:"projectID"`
	Blacklist bool   `json:"blacklist"`
}

type ListPlaylistIDsResponse struct {
	IDs []string `json:"IDs"`
}

type ListVideoIDsRequest struct {
	ProjectID string `json:"projectID"`
	Blacklist bool   `json:"blacklist"`
}

type ListVideoIDsResponse struct {
	IDs []string `json:"IDs"`
}
