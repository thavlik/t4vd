package server

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) handleGetChannelAvatar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusBadRequest)
				return errors.New("bad method")
			}
			channelID := r.URL.Query().Get("c")
			if channelID == "" {
				return errors.New("missing channelID in query")
			}
			noDownload := r.URL.Query().Get("nodownload") == "1"
			if noDownload {
				if err := s.scheduleChannelQuery(channelID); err != nil {
					return errors.Wrap(err, "scheduler.Add")
				}
				return nil
			}
			w.Header().Set("Content-Type", "image/jpeg")
			if err := s.thumbCache.Get(
				r.Context(),
				channelID,
				w,
			); err == api.ErrNotCached {
				if err := s.scheduleChannelQuery(channelID); err != nil {
					return errors.Wrap(err, "scheduler.Add")
				}
				w.WriteHeader(http.StatusNotFound)
				return nil
			} else if err != nil {
				return errors.Wrap(err, "cache.Get")
			}
			return nil
		}(); err != nil {
			s.log.Error("get channel avatar handler error", zap.Error(err))
		}
	}
}
