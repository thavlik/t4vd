package definitions

type SlideShow interface {
	GetRandomMarker(GetRandomMarker) Marker
	GetRandomStack(GetRandomStack) Stack
}

type Void struct{}

type GetRandomMarker struct {
	ProjectID string `json:"projectID"`
}

type Marker struct {
	VideoID string `json:"videoID"`
	Time    int64  `json:"time"`
}

type GetRandomStack struct {
	ProjectID string `json:"projectID"`
	Size      int    `json:"size"`
}

type Stack struct {
	Markers []*Marker `json:"markers"`
}
