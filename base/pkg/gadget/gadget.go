package gadget

type Gadget interface {
	ProjectID() string
	ID() string
	Name() string
}

type Dataset struct {
	ProjectID        string `json:"projectID"`
	Name             string `json:"name"`
	RestrictToUserID string `json:"restrictToUserID"`
}

type Filter struct {
	ProjectID string `json:"projectID"`
	Name      string `json:"name"`
	Source    Gadget
}

type Tagger struct {
	ProjectID string `json:"projectID"`
	Name      string `json:"name"`
	Source    Gadget
}

type Cropper struct {
	ProjectID string `json:"projectID"`
	Name      string `json:"name"`
}
