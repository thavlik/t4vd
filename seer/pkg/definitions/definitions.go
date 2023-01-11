package definitions

type Seer interface {
	GetChannelDetails(GetChannelDetailsRequest) GetChannelDetailsResponse
	GetPlaylistDetails(GetPlaylistDetailsRequest) GetPlaylistDetailsResponse
	GetVideoDetails(GetVideoDetailsRequest) GetVideoDetailsResponse
	GetBulkVideosDetails(GetBulkVideosDetailsRequest) GetBulkVideosDetailsResponse
	GetBulkPlaylistsDetails(GetBulkPlaylistsDetailsRequest) GetBulkPlaylistsDetailsResponse
	GetBulkChannelsDetails(GetBulkChannelsDetailsRequest) GetBulkChannelsDetailsResponse

	GetChannelVideoIDs(GetChannelVideoIDsRequest) GetChannelVideoIDsResponse
	GetPlaylistVideoIDs(GetPlaylistVideoIDsRequest) GetPlaylistVideoIDsResponse
	PurgeVideo(PurgeVideo) Void
	ListCache(ListCacheRequest) ListCacheResponse
	ScheduleVideoDownload(ScheduleVideoDownload) Void
	BulkScheduleVideoDownloads(BulkScheduleVideoDownloads) Void
	ListVideoDownloads(Void) VideoDownloads
	CancelVideoDownload(CancelVideoDownload) Void
}

type Void struct{}

type VideoDownloads struct {
	VideoIDs []string `json:"videoIDs"`
}

type CancelVideoDownload struct {
	VideoID string `json:"videoID"`
}

type ScheduleVideoDownload struct {
	VideoID string `json:"videoID"`
}

type BulkScheduleVideoDownloads struct {
	VideoIDs []string `json:"videoIDs"`
}

type GetChannelDetailsRequest struct {
	Input string `json:"input"`
	Force bool   `json:"force"`
}

type ChannelDetails struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Subs   string `json:"subs"`
}

type GetChannelDetailsResponse struct {
	Details ChannelDetails `json:"details"`
}

type GetPlaylistDetailsRequest struct {
	Input string `json:"input"`
	Force bool   `json:"force"`
}

type PlaylistDetails struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Channel   string `json:"channel"`
	ChannelID string `json:"channelID"`
	NumVideos int    `json:"numVideos"`
}

type GetPlaylistDetailsResponse struct {
	Details PlaylistDetails `json:"details"`
}

type GetVideoDetailsRequest struct {
	Input string `json:"input"`
	Force bool   `json:"force"`
}

type GetBulkVideosDetailsRequest struct {
	VideoIDs []string `json:"videoIDs"`
}

type GetBulkVideosDetailsResponse struct {
	Videos []*VideoDetails `json:"videos"`
}

type GetBulkPlaylistsDetailsRequest struct {
	PlaylistIDs []string `json:"playlistIDs"`
}

type GetBulkPlaylistsDetailsResponse struct {
	Playlists []*PlaylistDetails `json:"playlists"`
}

type GetBulkChannelsDetailsRequest struct {
	ChannelIDs []string `json:"channelIDs"`
}

type GetBulkChannelsDetailsResponse struct {
	Channels []*ChannelDetails `json:"channels"`
}

type VideoDetails struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	UploadDate  string `json:"uploadDate"`
	Uploader    string `json:"uploader"`
	UploaderID  string `json:"uploaderID"`
	Channel     string `json:"channel"`
	ChannelID   string `json:"channelID"`
	Duration    int64  `json:"duration"`
	ViewCount   int64  `json:"viewCount"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	FPS         int    `json:"fps"`
}

type GetVideoDetailsResponse struct {
	Details VideoDetails `json:"details"`
}

type GetChannelVideoIDsRequest struct {
	ID string `json:"id"`
}

type GetChannelVideoIDsResponse struct {
	VideoIDs []string `json:"videoIDs"`
}

type GetPlaylistVideoIDsRequest struct {
	ID string `json:"id"`
}

type GetPlaylistVideoIDsResponse struct {
	VideoIDs []string `json:"videoIDs"`
}

type PurgeVideo struct {
	ID string `json:"id"`
}

type ListCacheRequest struct {
	Marker string `json:"marker"`
}

type ListCacheResponse struct {
	VideoIDs    []string `json:"videoIDs"`
	IsTruncated bool     `json:"isTruncated"`
	NextMarker  string   `json:"nextMarker"`
}
