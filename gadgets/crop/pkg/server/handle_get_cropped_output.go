package server

import (
	"fmt"
	"image"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/gadget"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
	"go.uber.org/zap"
)

func handleGetCroppedOutput(
	labelStore labelstore.LabelStore,
	ref *gadget.DataRef,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			if r.Method != http.MethodGet {
				retCode = http.StatusMethodNotAllowed
				return errors.New("method not allowed")
			}
			id := r.URL.Query().Get("id")
			if id == "" {
				retCode = http.StatusBadRequest
				return errors.New("id is required")
			}
			gadgetName, channel, err := ref.Get(r.Context())
			if err == gadget.ErrNullDataRef {
				retCode = http.StatusNotFound
				return err
			} else if err != nil {
				return err
			}
			label, err := labelStore.Get(r.Context(), id)
			if err == labelstore.ErrNotFound {
				retCode = http.StatusNotFound
				return err
			} else if err != nil {
				return errors.Wrap(err, "failed to get label")
			}
			box, err := extractBbox(label.Payload)
			if err != nil {
				return errors.Wrap(err, "failed to extract bbox from label")
			}
			url := fmt.Sprintf(
				"%s/output/%s/x?%s",
				gadget.ResolveBaseURL(gadgetName),
				channel,
				r.URL.Query().Encode(),
			)
			req, err := http.NewRequest(
				http.MethodGet,
				url,
				nil,
			)
			req = req.WithContext(r.Context())
			if err != nil {
				return errors.Wrap(err, "failed to create request")
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return errors.Wrap(err, "failed to get input data")
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				w.WriteHeader(resp.StatusCode)
				body, _ := io.ReadAll(resp.Body)
				log.Warn(r.RequestURI, zap.Error(errors.Errorf(
					"%s: %s: %s",
					url,
					resp.Status,
					string(body),
				)))
				return nil
			}
			ct := resp.Header.Get("Content-Type")
			w.Header().Set("Content-Type", ct)
			switch ct {
			case "video/webm":
				marker, err := extractMarker(label.Payload)
				if err != nil {
					return errors.Wrap(err, "failed to extract marker from label")
				}
				return cropVideo(
					box,
					marker,
					resp.Body,
					w,
				)
			case "image/jpeg":
				fallthrough
			case "image/png":
				return cropImage(
					box,
					resp.Body,
					w,
				)
			default:
				retCode = http.StatusUnsupportedMediaType
				return errors.Errorf("unsupported media type: %s", ct)
			}
		}(); err != nil {
			gadget.HandlerError(r, w, retCode, err, log)
		}
	}
}

type Marker struct {
	Start time.Duration
	End   time.Duration
}

func extractMarker(
	payload map[string]interface{},
) (*Marker, error) {
	marker := &Marker{}
	v := false
	if start, ok := payload["t0"].(int64); ok {
		v = true
		marker.Start = time.Duration(start)
	}
	if end, ok := payload["t1"].(int64); ok {
		v = true
		marker.End = time.Duration(end)
	}
	if !v {
		return nil, nil
	}
	return marker, nil
}

func extractBbox(
	payload map[string]interface{},
) (image.Rectangle, error) {
	var x0, y0, x1, y1 int
	if v, ok := payload["x0"].(float64); ok {
		x0 = int(v)
	} else {
		return image.Rectangle{}, errors.New("missing x0")
	}
	if v, ok := payload["y0"].(float64); ok {
		y0 = int(v)
	} else {
		return image.Rectangle{}, errors.New("missing y0")
	}
	if v, ok := payload["x1"].(float64); ok {
		x1 = int(v)
	} else {
		return image.Rectangle{}, errors.New("missing x1")
	}
	if v, ok := payload["y1"].(float64); ok {
		y1 = int(v)
	} else {
		return image.Rectangle{}, errors.New("missing y1")
	}
	return image.Rect(
		x0,
		y0,
		x1,
		y1,
	), nil
}
