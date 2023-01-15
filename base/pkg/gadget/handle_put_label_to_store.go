package gadget

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/thavlik/t4vd/filter/pkg/api"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// HandlePutOutputLabel handler for inserting a label into the
// output channel. The frontend invokes this to insert
// labels into the host gadget's underlying storage.
func HandlePutOutputLabel(
	labelStore labelstore.LabelStore,
	payloadValidator PayloadValidator,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			if r.Method != http.MethodPut {
				retCode = http.StatusMethodNotAllowed
				return errors.Errorf("invalid method: %s", r.Method)
			}
			if r.Header.Get("Content-Type") != "application/json" {
				retCode = http.StatusBadRequest
				return errors.Errorf("invalid content type: %s", r.Header.Get("Content-Type"))
			}
			var label api.Label
			if err := json.NewDecoder(r.Body).Decode(&label); err != nil {
				retCode = http.StatusBadRequest
				return errors.Wrap(err, "json.Decode")
			}
			if label.ID == "" {
				retCode = http.StatusBadRequest
				return errors.New("label id must not be empty")
			}
			if label.ProjectID == "" {
				retCode = http.StatusBadRequest
				return errors.New("label project id must not be empty")
			}
			if label.CreatorID == "" {
				retCode = http.StatusBadRequest
				return errors.New("label creator id must not be empty")
			}
			if label.Parent == nil || label.Parent.ID == "" {
				retCode = http.StatusBadRequest
				return errors.New("label parent id must not be empty")
			}
			if payloadValidator != nil {
				if err := payloadValidator(label.Payload); err != nil {
					retCode = http.StatusBadRequest
					return errors.Wrap(err, "payload validation failed")
				}
			}
			now := time.Now()
			label.Created = &now
			if err := labelStore.Insert(&label); err != nil {
				return errors.Wrap(err, "failed to insert label")
			}
			log.Debug("inserted label",
				zap.String("projectID", label.ProjectID))
			return nil
		}(); err != nil {
			HandlerError(r, w, retCode, err, log)
		}
	}
}
