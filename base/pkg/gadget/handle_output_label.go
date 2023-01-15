package gadget

import (
	"net/http"

	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
	"go.uber.org/zap"
)

type PayloadValidator func(map[string]interface{}) error

// HandleOutputLabel samples labels from the output channel.
// This is gadget-specific as it requires querying the
// gadget's underlying storage.
func HandleOutputLabel(
	labelStore labelstore.LabelStore,
	ref *DataRef,
	payloadValidator PayloadValidator,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			switch r.Method {
			case http.MethodGet:
				HandleGetOutputLabel(labelStore, ref, log)(w, r)
				return nil
			case http.MethodPut:
				HandlePutOutputLabel(labelStore, payloadValidator, log)(w, r)
				return nil
			case http.MethodDelete:
				HandleDeleteOutputLabel(labelStore, log)(w, r)
				return nil
			default:
				retCode = http.StatusMethodNotAllowed
				return base.InvalidMethod(r.Method)
			}
		}(); err != nil {
			HandlerError(r, w, retCode, err, log)
		}
	}
}
