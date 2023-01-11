package server

import (
	"net/http"
	"time"

	"github.com/thavlik/t4vd/base/pkg/iam"
	"go.uber.org/zap"
)

type Subscription struct {
	userID    string
	projectID string
	stop      chan struct{}
	ch        chan []byte
}

func (s *Server) subscribe(userID, projectID string) *Subscription {
	sub := &Subscription{
		userID:    userID,
		projectID: projectID,
		stop:      make(chan struct{}, 1),
		ch:        make(chan []byte, 32),
	}
	s.subsL.Lock()
	s.subs[projectID] = append(s.subs[projectID], sub)
	s.subsL.Unlock()
	return sub
}

func (s *Server) unsubscribe(sub *Subscription) {
	sub.stop <- struct{}{}
	s.subsL.Lock()
	defer s.subsL.Unlock()
	subs := s.subs[sub.projectID]
	var newSubs []*Subscription
	for _, s := range subs {
		if sub == s {
			continue
		}
		newSubs = append(newSubs, sub)
	}
	if len(newSubs) == 0 {
		delete(s.subs, sub.projectID)
	} else {
		s.subs[sub.projectID] = newSubs
	}
}

func (s *Server) handleServerSentEvents() http.HandlerFunc {
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
			if err := s.ProjectAccess(r.Context(), userID, projectID); err != nil {

				w.WriteHeader(http.StatusForbidden)
				return nil
			}
			reqLog := s.log.With(
				zap.String("userID", userID),
				zap.String("projectID", projectID))
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
