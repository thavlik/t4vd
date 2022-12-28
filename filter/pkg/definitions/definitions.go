package definitions

type Filter interface {
	GetStack(GetStack) Stack
	Classify(Classify) Void
}

type GetStack struct {
	ProjectID string `json:"projectID"`
}

type Marker struct {
	VideoID string `json:"videoID"`
	Time    int64  `json:"time"`
}

type Stack struct {
	Markers []*Marker `json:"markers"`
}

type Classify struct {
	ProjectID string `json:"projectID"`
	Marker    Marker `json:"marker"`
	Label     int64  `json:"label"`
}

type Void struct{}
