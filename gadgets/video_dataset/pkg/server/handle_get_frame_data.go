package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/gadget"
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
			videoID := r.URL.Query().Get("v")
			if videoID == "" {
				retCode = http.StatusBadRequest
				return errors.New("invalid video id")
			}
			t := r.URL.Query().Get("t")
			if t == "" {
				retCode = http.StatusBadRequest
				return errors.New("missing time offset parameter")
			}
			timeOffset, err := strconv.ParseInt(t, 10, 64)
			if err != nil {
				retCode = http.StatusBadRequest
				return errors.Wrap(err, "invalid time offset")
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
