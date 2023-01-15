package server

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/gadget"
	"github.com/thavlik/t4vd/filter/pkg/api"
	seer "github.com/thavlik/t4vd/seer/pkg/api"
	"go.uber.org/zap"
)

func handleGetVideoMeta(
	seerClient seer.Seer,
	gadgetID string,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			id := r.URL.Query().Get("id")
			if id == "" {
				retCode = http.StatusBadRequest
				return errors.New("missing id parameter")
			}
			return json.NewEncoder(w).Encode(&api.Label{
				ID:        id,
				GadgetID:  gadgetID,
				CreatorID: "",  // TODO: resolve the userID of the user that submitted the video
				Created:   nil, // TODO: resolve the time the video was submitted
				Payload: map[string]interface{}{
					"videoID": id,
				},
			})
		}(); err != nil {
			gadget.HandlerError(r, w, retCode, err, log)
		}
	}
}
