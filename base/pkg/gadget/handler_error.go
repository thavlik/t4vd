package gadget

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func HandlerError(
	r *http.Request,
	w http.ResponseWriter,
	retCode int,
	err error,
	log *zap.Logger,
) {
	log.Error(r.RequestURI, zap.Error(err))
	w.WriteHeader(retCode)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
