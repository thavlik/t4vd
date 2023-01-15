package gadget

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/api"
	"go.uber.org/zap"
)

// handleGetOutputDataFromRef returns the data associated with input label.
// This is a proxy method that calls the input gadget's /output/x
// without transforming the data in any way.
func HandleGetOutputDataFromRef(
	gadgetID string,
	ref *DataRef,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			defer r.Body.Close()
			if r.Method != http.MethodPost {
				retCode = http.StatusMethodNotAllowed
				return errors.Errorf("method not allowed: %s", r.Method)
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
			body, err := json.Marshal(label.Parent)
			if err != nil {
				retCode = http.StatusInternalServerError
				return errors.Wrap(err, "json.Marshal")
			}
			gadgetName, channel, err := ref.Get(r.Context())
			if err != nil {
				return err
			}
			url := fmt.Sprintf(
				"%s/output/%s/x",
				ResolveBaseURL(gadgetName),
				channel,
			)
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
				retCode = resp.StatusCode
				body, _ := io.ReadAll(resp.Body)
				return errors.Errorf(
					"%s: %s: %s",
					url,
					resp.Status,
					string(body),
				)
			}
			if v := resp.Header.Get("Content-Type"); v != "" {
				w.Header().Set("Content-Type", v)
			}
			if v := resp.Header.Get("Content-Length"); v != "" {
				w.Header().Set("Content-Length", v)
			}
			if _, err := io.Copy(w, resp.Body); err != nil {
				return errors.Wrap(err, "failed to copy response")
			}
			return nil
		}(); err != nil {
			HandlerError(r, w, retCode, err, log)
		}
	}
}
