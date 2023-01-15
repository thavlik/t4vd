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

func extractMarkerID(id string) (string, int64, error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return "", 0, errors.New("invalid id parameter")
	}
	videoID := parts[0]
	if videoID == "" {
		return "", 0, errors.New("invalid video id")
	}
	timeOffset, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", 0, errors.Wrap(err, "invalid time offset")
	}
	return videoID, timeOffset, nil
}

func handleGetFrameMeta(
	seerClient seer.Seer,
	gadgetID string,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			videoID, timeOffset, err := extractMarkerID(r.URL.Query().Get("id"))
			if err != nil {
				retCode = http.StatusBadRequest
				return errors.Wrap(err, "invalid id parameter")
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
