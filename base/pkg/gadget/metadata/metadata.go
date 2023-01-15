package metadata

type Channel struct {
	Name string `json:"name"`
}

type Metadata struct {
	Name         string     `json:"name"`
	MaxBatchSize int        `json:"maxBatchSize"`
	Inputs       []*Channel `json:"inputs"`
	Outputs      []*Channel `json:"outputs"`
}
