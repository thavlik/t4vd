package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/gadget"
	"github.com/thavlik/t4vd/filter/pkg/api"
	slideshow "github.com/thavlik/t4vd/slideshow/pkg/api"
	"go.uber.org/zap"
)

func handleSampleRandomFrames(
	compilerOpts base.ServiceOptions,
	slideshowOpts base.ServiceOptions,
	gadgetID string,
	projectID string,
	maxBatchSize int,
	log *zap.Logger,
) http.HandlerFunc {
	//compilerClient := compiler.NewCompilerClientFromOptions(compilerOpts)
	slideshowClient := slideshow.NewSlideShowClientFromOptions(slideshowOpts)
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			batchSize, err := gadget.ExtractBatchSize(r.URL.Query(), maxBatchSize)
			if err != nil {
				retCode = http.StatusBadRequest
				return err
			}
			labels := make([]*api.Label, batchSize)
			for i := 0; i < batchSize; i++ {
				marker, err := slideshowClient.GetRandomMarker(
					r.Context(),
					slideshow.GetRandomMarker{
						ProjectID: projectID,
					})
				if err != nil {
					return err
				}
				labels = append(labels, &api.Label{
					ID: fmt.Sprintf(
						"%s:%d",
						marker.VideoID,
						marker.Time,
					),
					GadgetID: gadgetID,
					Payload: map[string]interface{}{
						"timeOffset": marker.Time,
					},
					Parent: &api.Label{
						ID:       marker.VideoID,
						GadgetID: gadgetID,
						Payload: map[string]interface{}{
							"videoID": marker.VideoID,
						},
					},
				})
			}
			w.Header().Set("Content-Type", "application/json")
			return json.NewDecoder(r.Body).Decode(&labels)
		}(); err != nil {
			gadget.HandlerError(r, w, retCode, err, log)
		}
	}
}
