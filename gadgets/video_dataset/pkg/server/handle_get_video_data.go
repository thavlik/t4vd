package server

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/gadget"
	"github.com/thavlik/t4vd/filter/pkg/api"
	"go.uber.org/zap"
)

func handleGetVideoData(
	bucketName string,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			defer r.Body.Close()
			if r.Method != http.MethodPost {
				retCode = http.StatusMethodNotAllowed
				return errors.New("method not allowed")
			}
			if r.Header.Get("Content-Type") != "application/json" {
				retCode = http.StatusBadRequest
				return errors.New("Content-Type is not application/json")
			}
			var label api.Label
			if err := json.NewDecoder(r.Body).Decode(&label); err != nil {
				retCode = http.StatusBadRequest
				return errors.Wrap(err, "json.Decode")
			}
			if label.Payload == nil {
				retCode = http.StatusBadRequest
				return errors.New("label.Payload is nil")
			}
			videoID, ok := label.Payload["videoID"].(string)
			if !ok {
				retCode = http.StatusBadRequest
				return errors.New("label.Payload[\"videoID\"] is not a string")
			}
			// TODO: get video from bucket
			_ = videoID
			w.WriteHeader(http.StatusNotImplemented)
			return nil
		}(); err != nil {
			gadget.HandlerError(r, w, retCode, err, log)
		}
	}
}
