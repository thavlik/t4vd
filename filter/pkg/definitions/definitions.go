package definitions

type Filter interface {
	Classify(Classify) Void
}

type Marker struct {
	VideoID string `json:"videoID"`
	Time    int64  `json:"time"`
}

type Classify struct {
	ProjectID string `json:"projectID"`
	Marker    Marker `json:"marker"`
	Label     int64  `json:"label"`
}

type Void struct{}
