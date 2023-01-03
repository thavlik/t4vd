package server

import (
	"net/http"
	"time"

	"github.com/thavlik/t4vd/base/pkg/iam"
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
			// TODO: make sure user has access to project
			done := r.Context().Done()
			sub := s.subscribe(userID, projectID)
			defer s.unsubscribe(sub)
			for {
				select {
				case <-done:
					return r.Context().Err()
				case <-time.After(5 * time.Second):
					w.Write([]byte("{}\n"))
				case msg := <-sub.ch:
					w.Write(msg)
					w.Write([]byte("\n"))
				}
			}
		})
}
