package gadget

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"

	"go.uber.org/zap"
)

func HandleSampleOutputLabelsFromStore(
	labelStore labelstore.LabelStore,
	maxBatchSize int,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			if r.Method != http.MethodGet {
				retCode = http.StatusMethodNotAllowed
				return errors.Errorf("invalid method: %s", r.Method)
			}
			batchSize, err := ExtractBatchSize(
				r.URL.Query(),
				maxBatchSize,
			)
			if err != nil {
				retCode = http.StatusBadRequest
				return errors.Wrap(err, "failed to extract batch size from query")
			}
			projectID := r.URL.Query().Get("p")
			if projectID == "" {
				retCode = http.StatusBadRequest
				return errors.New("missing project ID")
			}
			labels, err := labelStore.Sample(
				r.Context(),
				&labelstore.SampleInput{
					ProjectID: projectID,
					BatchSize: batchSize,
				},
			)
			if err != nil {
				return err
			}
			log.Debug("sampled output labels", zap.Int("count", len(labels)))
			w.Header().Set("Content-Type", "application/json")
			return json.NewEncoder(w).Encode(labels)
		}(); err != nil {
			HandlerError(r, w, retCode, err, log)
		}
	}
}
