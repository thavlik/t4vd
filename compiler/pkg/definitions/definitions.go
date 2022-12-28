package definitions

type Compiler interface {
	Compile(Compile) Void
	GetDataset(GetDatasetRequest) Dataset
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

type Video struct {
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

type Dataset struct {
	ID        string   `json:"id"`
	Timestamp int64    `json:"timestamp"`
	Complete  bool     `json:"complete"`
	Videos    []*Video `json:"videos"`
}
