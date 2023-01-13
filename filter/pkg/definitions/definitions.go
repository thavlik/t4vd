package definitions

type Filter interface {
	Classify(Label) Label
	Sample(SampleRequest) SampleResponse
}

type SampleRequest struct {
	ProjectID string `json:"projectID"`
	BatchSize int    `json:"batchSize"`
}

type SampleResponse struct {
	Labels []*Label `json:"labels"`
}

// Marker refers to a frame in a video
type Marker struct {
	VideoID   string `json:"videoID"`   // ID of video
	Timestamp int64  `json:"timestamp"` // time in video
}

// Label is a classification of a frame in a video.
// 'Tags' is temporary and will be replaced with a more
// robust system capable of handling multiple types of
// labeled data, such as XY coordinates, bounding boxes,
// etc.
type Label struct {
	ID          string   `json:"id"`          // unique ID for the label
	Timestamp   int64    `json:"timestamp"`   // time of classification
	ProjectID   string   `json:"projectID"`   // project ID
	Marker      Marker   `json:"marker"`      // video and timestamp in video of frame
	SubmitterID string   `json:"submitterID"` // ID of user who submitted label
	ParentID    string   `json:"parentID"`    // ID of parent label
	Tags        []string `json:"tags"`        // tags to classify frame with
}
