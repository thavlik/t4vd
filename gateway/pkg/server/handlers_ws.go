package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/iam"
	"go.uber.org/zap"
)

func (s *Server) handleSubscribe() websockMessageHandler {
	return func(
		ctx context.Context,
		userID string,
		msg map[string]interface{},
		c *websocket.Conn,
	) error {
		projectID, ok := msg["projectID"].(string)
		if !ok {
			return errors.New("invalid projectID")
		}
		if err := s.ProjectAccess(ctx, userID, projectID); err != nil {
			return errors.Wrap(err, "ProjectAccess")
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case s.wsSubsL <- struct{}{}:
			defer func() { <-s.wsSubsL }()
		}
		if subs, ok := s.wsSubs[c]; ok {
			if unsubAll, ok := msg["unsubscribeAll"].(bool); ok && unsubAll {
				// unsubscribe from everything else
				var newSubs []*Subscription
				for _, sub := range subs {
					if sub.projectID != projectID {
						s.unsubscribe(sub)
					} else {
						newSubs = append(newSubs, sub)
					}
				}
				if len(newSubs) > 1 {
					panic(fmt.Sprintf("sanity check failed: len(newSubs) > 1: %d", len(newSubs)))
				}
				s.wsSubs[c] = newSubs
				subs = newSubs
			}
			for _, sub := range subs {
				if sub.projectID == projectID {
					// already subscribed
					return nil
				}
			}
		}
		sub := s.subscribe(userID, projectID)
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case <-sub.stop:
					return
				case msg := <-sub.ch:
					if err := c.WriteMessage(0, msg); err != nil {
						s.log.Warn("WriteMessage", zap.Error(err))
					}
				}
			}
		}()
		s.wsSubs[c] = append(s.wsSubs[c], sub)
		return nil
	}
}

func (s *Server) handleUnsubscribe() websockMessageHandler {
	return func(
		ctx context.Context,
		userID string,
		msg map[string]interface{},
		c *websocket.Conn,
	) error {
		projectID, ok := msg["projectID"].(string)
		if !ok {
			return errors.New("invalid projectID")
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case s.wsSubsL <- struct{}{}:
			defer func() { <-s.wsSubsL }()
		}
		subs, ok := s.wsSubs[c]
		if !ok {
			return errors.New("no subscriptions")
		}
		for i, sub := range subs {
			if sub.projectID == projectID {
				s.unsubscribe(sub)
				subs = append(subs[:i], subs[i+1:]...)
				break
			}
		}
		s.wsSubs[c] = subs
		return nil
	}
}

func (s *Server) handlePing() websockMessageHandler {
	return func(
		ctx context.Context,
		userID string,
		msg map[string]interface{},
		c *websocket.Conn,
	) error {
		return c.WriteJSON(map[string]interface{}{
			"type": "pong",
		})
	}
}

func (s *Server) handleWebSock() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodGet,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			upgrader := websocket.Upgrader{}
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return fmt.Errorf("upgrade: %v", err)
			}
			ctx := r.Context()
			defer c.Close()
			defer func() {
				// clean up subscriptions
				s.wsSubsL <- struct{}{}
				defer func() { <-s.wsSubsL }()
				subs, ok := s.wsSubs[c]
				if !ok {
					return
				}
				for _, sub := range subs {
					s.unsubscribe(sub)
				}
				delete(s.wsSubs, c)
			}()
			messages := make(chan []byte, 32)
			stopped := make(chan error)
			go func() {
				defer close(messages)
				stopped <- func() error {
					for {
						_, buf, err := c.ReadMessage()
						if err != nil {
							return err
						}
						select {
						case <-ctx.Done():
							return ctx.Err()
						case messages <- buf:
							continue
						}
					}
				}()
			}()
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case err := <-stopped:
					return err
				case msg := <-messages:
					obj := make(map[string]interface{})
					if err := json.Unmarshal(msg, &obj); err != nil {
						return errors.Wrap(err, "unmarshal")
					}
					ty, ok := obj["type"].(string)
					if !ok {
						return errors.New("invalid message")
					}
					handler, ok := s.wsHandlers[ty]
					if !ok {
						return errors.New("unrecognized message type")
					}
					if err := handler(
						ctx,
						userID,
						obj,
						c,
					); err != nil {
						return errors.Wrap(err, ty)
					}
				}
			}
		},
	)
}
