package gadget

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"go.uber.org/zap"
)

// HandleGetOutputLabelFromRef a handler that retrieves
// output labels from the referenced gadget
func HandleGetOutputLabelFromRef(
	ref *DataRef,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			if r.Method != http.MethodPost {
				retCode = http.StatusMethodNotAllowed
				return base.InvalidMethod(r.Method)
			}
			label, err := GetOutputLabelFromRef(
				r.Context(),
				nil,
				ref,
				log,
			)
			if err == ErrLabelNotFound {
				retCode = http.StatusNotFound
				return err
			} else if err != nil {
				return errors.Wrap(err, "failed to get label")
			}
			w.Header().Set("Content-Type", "application/json")
			return json.NewEncoder(w).Encode(label)
		}(); err != nil {
			HandlerError(r, w, retCode, err, log)
		}
	}
}
