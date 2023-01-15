package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/gadget"
	"github.com/thavlik/t4vd/filter/pkg/api"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
	"go.uber.org/zap"
)

func handleGetCroppedOutput(
	labelStore labelstore.LabelStore,
	gadgetID string,
	ref *gadget.DataRef,
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
			if label.GadgetID != gadgetID {
				retCode = http.StatusBadRequest
				return errors.Errorf(
					"mismatched gadgetID: %s != %s",
					label.GadgetID,
					gadgetID,
				)
			}
			if label.Parent == nil {
				retCode = http.StatusBadRequest
				return errors.New("label.Parent is nil")
			}
			gadgetName, channel, err := ref.Get(r.Context())
			if err != nil {
				return err
			}
			url := fmt.Sprintf(
				"%s/output/%s/x",
				gadget.ResolveBaseURL(gadgetName),
				channel,
			)
			body, err := json.Marshal(label.Parent)
			if err != nil {
				return errors.Wrap(err, "failed to marshal label parent")
			}
			req, err := http.NewRequest(
				http.MethodPost,
				url,
				bytes.NewReader(body),
			)
			req.Header.Set("Content-Type", "application/json")
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
			if v := resp.Header.Get("Content-Length"); v != "" {
				w.Header().Set("Content-Length", v)
			}
			ct := resp.Header.Get("Content-Type")
			w.Header().Set("Content-Type", ct)
			switch ct {
			case "video/webm":
				box, err := extractBbox(label.Payload)
				if err != nil {
					retCode = http.StatusBadRequest
					return errors.Wrap(err, "failed to extract bbox from label")
				}
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
				box, err := extractBbox(label.Payload)
				if err != nil {
					retCode = http.StatusBadRequest
					return errors.Wrap(err, "failed to extract bbox from label")
				} else if box == nil {
					retCode = http.StatusBadRequest
					return errors.New("no bbox")
				}
				return cropImage(
					*box,
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
) (*image.Rectangle, error) {
	var x0, y0, x1, y1 int
	if v, ok := payload["x0"].(float64); ok {
		x0 = int(v)
	} else {
		return nil, nil
	}
	if v, ok := payload["y0"].(float64); ok {
		y0 = int(v)
	} else {
		return nil, nil
	}
	if v, ok := payload["x1"].(float64); ok {
		x1 = int(v)
	} else {
		return nil, nil
	}
	if v, ok := payload["y1"].(float64); ok {
		y1 = int(v)
	} else {
		return nil, nil
	}
	r := image.Rect(
		x0,
		y0,
		x1,
		y1,
	)
	return &r, nil
}
