package server

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/vidcache"
	"go.uber.org/zap"
)

func (s *Server) handleGetVideo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
			videoID := r.URL.Query().Get("v")
			if videoID == "" {
				return errors.New("missing videoID from query")
			}
			noDownload := r.URL.Query().Get("nodownload") == "1"
			log := s.log.With(zap.String("videoID", videoID))
			if noDownload {
				// just make sure it's cached by adding
				// the videoID to the scheduler
				if err := s.dlSched.Add(videoID); err != nil {
					return errors.Wrap(err, "scheduler.Add")
				}
				log.Debug("added video to scheduler")
				return nil
			}
			w.Header().Set("Content-Type", "video/webm")
			if err := s.vidCache.Get(
				r.Context(),
				videoID,
				w,
			); err == vidcache.ErrVideoNotCached {
				if err := s.dlSched.Add(videoID); err != nil {
					return errors.Wrap(err, "scheduler.Add")
				}
				log.Debug("added video to scheduler")
				w.WriteHeader(http.StatusNotFound)
				return nil
			} else if err != nil {
				return errors.Wrap(err, "cache.Get")
			}
			return nil
		}(); err != nil {
			s.log.Error("get video handler error", zap.Error(err))
		}
	}
}
