package gadget

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
	"go.uber.org/zap"
)

// HandleGetOutputLabel handler for retrieving a label from the
// output channel. The frontend invokes this to retrieve
// a specific label's metadata from the host gadget's underlying
// storage given its id.
func HandleGetOutputLabel(
	labelStore labelstore.LabelStore,
	ref *DataRef,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			if r.Method != http.MethodGet {
				retCode = http.StatusMethodNotAllowed
				return errors.Errorf("invalid method: %s", r.Method)
			}
			id := r.URL.Query().Get("id")
			if id == "" {
				retCode = http.StatusBadRequest
				return errors.New("id is required")
			}
			label, err := labelStore.Get(r.Context(), id)
			if err == labelstore.ErrNotFound {
				retCode = http.StatusNotFound
				return err
			} else if err != nil {
				return errors.Wrap(err, "labelstore.Get")
			}
			if label.Parent != nil {
				// reference the input gadget to resolve the parent label
				if label.Parent, err = GetOutputLabelFromRef(
					r.Context(),
					label.Parent,
					ref,
					log,
				); err != nil {
					return errors.Wrap(err, "failed to resolve parent label")
				}
			}
			w.Header().Set("Content-Type", "application/json")
			return json.NewEncoder(w).Encode(label)
		}(); err != nil {
			HandlerError(r, w, retCode, err, log)
		}
	}
}
