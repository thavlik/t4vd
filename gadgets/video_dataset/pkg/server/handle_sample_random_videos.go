package server

import (
	"encoding/json"
	"net/http"

	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/filter/pkg/api"
	"go.uber.org/zap"
)

func handleSampleRandomVideos(
	// TODO: implement a Sample method on the compiler
	compilerOpts base.ServiceOptions,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videoID := ""
		if err := func() error {
			labels := []*api.Label{{
				ID: videoID,
				Payload: map[string]interface{}{
					"videoID": videoID,
				},
			}}
			w.Header().Set("Content-Type", "application/json")
			return json.NewDecoder(r.Body).Decode(&labels)
		}(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(r.RequestURI, zap.Error(err))
		}
	}
}
