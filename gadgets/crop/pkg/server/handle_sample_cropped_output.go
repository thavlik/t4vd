package server

import (
	"net/http"

	"go.uber.org/zap"
)

func handleSampleCroppedOutput(
	maxBatchSize int,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
