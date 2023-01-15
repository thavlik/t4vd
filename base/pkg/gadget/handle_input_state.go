package gadget

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func HandleInputState(
	channels map[string]*DataRef,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			channel := mux.Vars(r)["channel"]
			ref, ok := channels[channel]
			if !ok {
				retCode = http.StatusNotFound
				return nil
			}
			switch r.Method {
			case http.MethodGet:
				gadgetName, channelName, err := ref.Get(r.Context())
				if err != nil {
					return err
				}
				w.Header().Set("Content-Type", "application/json")
				return json.NewEncoder(w).Encode(map[string]string{
					"gadget":  gadgetName,
					"channel": channelName,
				})
			case http.MethodPut:
				var req struct {
					Gadget  string `json:"gadget"`
					Channel string `json:"channel"`
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					retCode = http.StatusBadRequest
					return err
				}
				if err := ref.Set(
					r.Context(),
					req.Gadget,
					req.Channel,
				); err != nil {
					return err
				}
				return nil
			default:
				retCode = http.StatusMethodNotAllowed
				return errors.Errorf("invalid method: %s", r.Method)
			}
		}(); err != nil {
			HandlerError(r, w, retCode, err, log)
		}
	}
}
