package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/gadget"
	"github.com/thavlik/t4vd/filter/pkg/api"
	slideshow "github.com/thavlik/t4vd/slideshow/pkg/api"
	"go.uber.org/zap"
)

func handleGetFrameData(
	slideshowOpts base.ServiceOptions,
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
			timeOffset, ok := label.Payload["timeOffset"].(int64)
			if !ok {
				retCode = http.StatusBadRequest
				return errors.New("label.Payload[\"timeOffset\"] is not an int64")
			}
			if label.Parent == nil {
				retCode = http.StatusBadRequest
				return errors.New("label.Parent is nil")
			}
			if label.Parent.Payload == nil {
				retCode = http.StatusBadRequest
				return errors.New("label.Parent.Payload is nil")
			}
			videoID, ok := label.Parent.Payload["videoID"].(string)
			if !ok {
				retCode = http.StatusBadRequest
				return errors.New("label.Parent.Payload[\"videoID\"] is not a string")
			}
			w.Header().Set("Content-Type", "image/jpeg")
			if err := slideshow.GetFrame(
				r.Context(),
				slideshowOpts,
				videoID,
				time.Duration(timeOffset),
				w,
			); err != nil {
				return errors.Wrap(err, "slideshow.GetFrame")
			}
			return nil
		}(); err != nil {
			gadget.HandlerError(r, w, retCode, err, log)
		}
	}
}
