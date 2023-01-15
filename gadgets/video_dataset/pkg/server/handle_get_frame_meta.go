package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/gadget"
	"github.com/thavlik/t4vd/filter/pkg/api"
	seer "github.com/thavlik/t4vd/seer/pkg/api"

	"go.uber.org/zap"
)

func handleGetFrameMeta(
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
			parts := strings.Split(id, ":")
			if len(parts) != 2 {
				retCode = http.StatusBadRequest
				return errors.New("invalid id parameter")
			}
			videoID := parts[0]
			if videoID == "" {
				retCode = http.StatusBadRequest
				return errors.New("invalid video id")
			}
			timeOffset, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				retCode = http.StatusBadRequest
				return errors.Wrap(err, "invalid time offset")
			}
			return json.NewEncoder(w).Encode(&api.Label{
				GadgetID: gadgetID,
				Payload: map[string]interface{}{
					"timeOffset": timeOffset,
				},
				Parent: &api.Label{
					GadgetID:  gadgetID,
					CreatorID: "",  // TODO: resolve the userID of the user that submitted the video
					Created:   nil, // TODO: resolve the time the video was submitted
					Payload: map[string]interface{}{
						"videoID": videoID,
					},
				},
			})
		}(); err != nil {
			gadget.HandlerError(r, w, retCode, err, log)
		}
	}
}
