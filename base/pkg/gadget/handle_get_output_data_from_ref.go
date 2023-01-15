package gadget

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// handleGetOutputDataFromRef returns the data associated with input label.
// This is a proxy method that calls the input gadget's /output/x
func HandleGetOutputDataFromRef(
	ref *DataRef,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			if r.Method != http.MethodGet {
				retCode = http.StatusMethodNotAllowed
				return errors.Errorf("method not allowed: %s", r.Method)
			}
			if !r.URL.Query().Has("id") {
				retCode = http.StatusBadRequest
				return errors.Errorf("id is required")
			}
			gadgetName, channel, err := ref.Get(r.Context())
			if err == ErrNullDataRef {
				retCode = http.StatusNotFound
				return err
			} else if err != nil {
				return err
			}
			url := fmt.Sprintf(
				"%s/output/%s/x?%s",
				ResolveBaseURL(gadgetName),
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
