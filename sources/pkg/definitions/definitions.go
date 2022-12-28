package definitions

type Sources interface {
	CreateProject(Project) Project
	DeleteProject(DeleteProject) Void
	ListProjects(ListProjectsRequest) ListProjectsResponse
	GetProject(GetProject) Project
	GetProjectByName(GetProjectByName) Project

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
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	Subs      string `json:"subs"`
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
	Title     string `json:"title"`
	Channel   string `json:"channel"`
	ChannelID string `json:"channelID"`
	NumVideos int    `json:"numVideos"`
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
	ID          string `json:"id"`
	Title       string `json:"title"`
	Blacklist   bool   `json:"blacklist"`
	Channel     string `json:"channel"`
	ChannelID   string `json:"channelID"`
	Description string `json:"description"`
	Duration    int64  `json:"duration"`
	FPS         int    `json:"fps"`
	Height      int    `json:"height"`
	Width       int    `json:"width"`
	Thumbnail   string `json:"thumbnail"`
	Uploader    string `json:"uploader"`
	UploaderID  string `json:"uploaderID"`
	ViewCount   int64  `json:"viewCount"`
	UploadDate  string `json:"uploadDate"`
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
