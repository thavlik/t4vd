package gadget

import (
	"encoding/json"
	"net/http"

	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/gadget/metadata"
	"go.uber.org/zap"
)

func HandleGetMetadata(
	metadata *metadata.Metadata,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			if r.Method != http.MethodGet {
				retCode = http.StatusMethodNotAllowed
				return base.InvalidMethod(r.Method)
			}
			return json.NewEncoder(w).Encode(metadata)
		}(); err != nil {
			HandlerError(r, w, retCode, err, log)
		}
	}
}
