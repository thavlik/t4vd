package definitions

type Compiler interface {
	Compile(Compile) Void
	GetDataset(GetDatasetRequest) Dataset
	ResolveProjectsForVideo(ResolveProjectsForVideoRequest) ResolveProjectsForVideoResponse
}

type ResolveProjectsForVideoRequest struct {
	VideoID string `json:"videoID"`
}

type ResolveProjectsForVideoResponse struct {
	ProjectIDs []string `json:"projectIDs"`
}

type Compile struct {
	ProjectID string `json:"projectID"`
	All       bool   `json:"all"`
}

type Void struct{}

type GetDatasetRequest struct {
	ID        string `json:"id,omitempty"`
	ProjectID string `json:"projectID"`
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

type VideoSource struct {
	Type        string `json:"type"`        // [playlist|channel]
	ID          string `json:"id"`          // playlistID or channelID
	SubmitterID string `json:"submitterID"` // userID of submitter
	Submitted   int64  `json:"submitted"`   // unix timestamp of submission
}

type Video struct {
	ID      string        `json:"id"`               // videoID
	Details *VideoDetails `json:"details"`          // details from youtube
	Source  *VideoSource  `json:"source,omitempty"` // nil if video is included directly
}

type Dataset struct {
	ID        string   `json:"id"`
	Timestamp int64    `json:"timestamp"`
	Complete  bool     `json:"complete"`
	Videos    []*Video `json:"videos"`
}
