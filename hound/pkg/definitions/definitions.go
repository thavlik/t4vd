package definitions

type Hound interface {
	ReportChannelVideo(ChannelVideo) Void
	ReportPlaylistVideo(PlaylistVideo) Void
	ReportVideoDetails(VideoDetails) Void
	ReportChannelDetails(ChannelDetails) Void
	ReportPlaylistDetails(PlaylistDetails) Void
	ReportVideoDownloadProgress(VideoDownloadProgress) Void
}

type Void struct{}

type ChannelVideo struct {
	ChannelID string       `json:"channelID"`
	NumVideos int          `json:"numVideos"`
	Video     VideoDetails `json:"video"`
}

type PlaylistVideo struct {
	PlaylistID string       `json:"playlistID"`
	NumVideos  int          `json:"numVideos"`
	Video      VideoDetails `json:"video"`
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

type ChannelDetails struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Subs   string `json:"subs"`
}

type PlaylistDetails struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Channel   string `json:"channel"`
	ChannelID string `json:"channelID"`
	NumVideos int    `json:"numVideos"`
}

type VideoDownloadProgress struct {
	ID      string  `json:"id"`
	Total   int64   `json:"total"`
	Rate    float64 `json:"rate"`
	Elapsed int64   `json:"elapsed"`
}
