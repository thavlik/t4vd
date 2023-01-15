package gadget

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"go.uber.org/zap"
)

func HandleSampleOutputLabelsFromRef(
	ref *DataRef,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			if r.Method != http.MethodGet {
				retCode = http.StatusMethodNotAllowed
				return base.InvalidMethod(r.Method)
			}
			gadgetName, channel, err := ref.Get(r.Context())
			if err != nil {
				return errors.Wrap(err, "failed to get data ref")
			}
			url := fmt.Sprintf(
				"%s/sample/output/%s/y?%s",
				ResolveBaseURL(gadgetName),
				channel,
				r.URL.Query().Encode(),
			)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				return errors.Wrap(err, "failed to create request")
			}
			req = req.WithContext(r.Context())
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return errors.Wrap(err, "failed to make request")
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				retCode = resp.StatusCode
				body, _ := io.ReadAll(resp.Body)
				return errors.Errorf("%s: %s: %s", url, resp.Status, body)
			}
			w.Header().Set("Content-Type", "application/json")
			if _, err := io.Copy(w, resp.Body); err != nil {
				return errors.Wrap(err, "failed to copy response")
			}
			return nil
		}(); err != nil {
			HandlerError(r, w, retCode, err, log)
		}
	}
}
