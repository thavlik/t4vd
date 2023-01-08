package server

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/iam"
	"go.uber.org/zap"
)

type Subscription struct {
	userID    string
	projectID string
	ch        chan []byte
}

func (s *Server) subscribe(userID, projectID string) *Subscription {
	sub := &Subscription{
		userID:    userID,
		projectID: projectID,
		ch:        make(chan []byte, 16),
	}
	s.subsL.Lock()
	s.subs[sub] = struct{}{}
	s.subsL.Unlock()
	return sub
}

func (s *Server) unsubscribe(sub *Subscription) {
	s.subsL.Lock()
	delete(s.subs, sub)
	s.subsL.Unlock()
}

func (s *Server) handleEvents() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodGet,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			projectID := r.URL.Query().Get("p")
			if projectID == "" {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
			reqLog := s.log.With(
				zap.String("userID", userID),
				zap.String("projectID", projectID))
			// make sure user has access to project
			if inGroup, err := s.iam.IsUserInGroup(r.Context(), userID, projectID); err != nil {
				return errors.Wrap(err, "iam.IsUserInGroup")
			} else if !inGroup {
				reqLog.Warn("user not in project group")
				w.WriteHeader(http.StatusForbidden)
				return nil
			}
			reqLog.Debug("subscribing to events")
			defer reqLog.Debug("unsubscribing from events")
			done := r.Context().Done()
			sub := s.subscribe(userID, projectID)
			defer s.unsubscribe(sub)
			for {
				select {
				case <-done:
					return r.Context().Err()
				case <-time.After(5 * time.Second):
					// send a keep-alive
					if _, err := w.Write([]byte("{}\n")); err != nil {
						return err
					}
				case msg := <-sub.ch:
					if _, err := w.Write(msg); err != nil {
						return err
					}
					if _, err := w.Write([]byte("\n")); err != nil {
						return err
					}
				}
			}
		})
}
