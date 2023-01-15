package server

import (
	"encoding/json"
	"net/http"

	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/gadget"
	compiler "github.com/thavlik/t4vd/compiler/pkg/api"
	"github.com/thavlik/t4vd/filter/pkg/api"
	"go.uber.org/zap"
)

func handleSampleRandomVideos(
	compilerOpts base.ServiceOptions,
	gadgetID string,
	projectID string,
	maxBatchSize int,
	log *zap.Logger,
) http.HandlerFunc {
	compilerClient := compiler.NewCompilerClientFromOptions(compilerOpts)
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			batchSize, err := gadget.ExtractBatchSize(r.URL.Query(), maxBatchSize)
			if err != nil {
				retCode = http.StatusBadRequest
				return err
			}
			videos, err := sampleRandomVideos(
				r.Context(),
				compilerClient,
				projectID,
				batchSize,
			)
			if err != nil {
				return err
			}
			labels := make([]*api.Label, batchSize)
			for i, video := range videos {
				labels[i] = videoToLabel(video, gadgetID)
			}
			w.Header().Set("Content-Type", "application/json")
			return json.NewDecoder(r.Body).Decode(&labels)
		}(); err != nil {
			gadget.HandlerError(r, w, retCode, err, log)
		}
	}
}

func videoToLabel(video *compiler.Video, gadgetID string) *api.Label {
	return &api.Label{
		GadgetID: gadgetID,
		Payload: map[string]interface{}{
			"videoID": video.Details.ID,
		},
	}
}
