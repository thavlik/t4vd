package server

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) handleGetVideoThumbnail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusBadRequest)
				return errors.New("bad method")
			}
			videoID := r.URL.Query().Get("v")
			if videoID == "" {
				return errors.New("missing videoID in query")
			}
			noDownload := r.URL.Query().Get("nodownload") == "1"
			if noDownload {
				if err := s.scheduleVideoQuery(videoID); err != nil {
					return err
				}
				return nil
			}
			w.Header().Set("Content-Type", "image/jpeg")
			if err := s.thumbCache.Get(
				r.Context(),
				videoID,
				w,
			); err == api.ErrNotCached {
				if err := s.scheduleVideoQuery(videoID); err != nil {
					return err
				}
				w.WriteHeader(http.StatusNotFound)
				return nil
			} else if err != nil {
				return errors.Wrap(err, "cache.Get")
			}
			return nil
		}(); err != nil {
			s.log.Error("get video thumbnail handler error", zap.Error(err))
		}
	}
}
