package gadget

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"

	"go.uber.org/zap"
)

// HandleDeleteOutputLabel deletes a label from the output channel.
// deleterID is required as it is used to track who deleted the label.
func HandleDeleteOutputLabel(
	labelStore labelstore.LabelStore,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			if r.Header.Get("Content-Type") != "application/json" {
				retCode = http.StatusMethodNotAllowed
				return errors.Errorf("invalid content type: %s", r.Header.Get("Content-Type"))
			}
			if r.Method != http.MethodDelete {
				retCode = http.StatusMethodNotAllowed
				return base.InvalidMethod(r.Method)
			}
			var resp struct {
				ID        string `json:"id"`
				DeleterID string `json:"deleterID"`
			}
			if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
				retCode = http.StatusBadRequest
				return errors.Wrap(err, "json.Decode")
			}
			if err := labelStore.Delete(
				&labelstore.DeleteInput{
					ID:        resp.ID,
					DeleterID: resp.DeleterID,
					Timestamp: time.Now(),
				},
			); err != nil {
				return errors.Wrap(err, "labelstore.Delete")
			}
			log.Debug("deleted label", zap.String("id", resp.ID))
			return nil
		}(); err != nil {
			HandlerError(r, w, retCode, err, log)
		}
	}
}
