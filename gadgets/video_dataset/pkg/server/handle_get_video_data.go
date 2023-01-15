package server

import (
	"net/http"

	"go.uber.org/zap"
)

func handleGetVideoData(
	bucketName string,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: download video directly from s3 bucket
	}
}
