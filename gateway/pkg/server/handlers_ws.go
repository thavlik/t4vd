package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/iam"
	"go.uber.org/zap"
)

const pingInterval = 10 * time.Second
const heartbeatTimeout = 20 * time.Second

func (s *Server) setupWebSockHandlers() {
	s.registerWebsockHandler("subscribe", s.handleSubscribe())
	s.registerWebsockHandler("unsubscribe", s.handleUnsubscribe())
	s.registerWebsockHandler("ping", s.handlePing())
}

func (s *Server) registerWebsockHandler(name string, handler websockMessageHandler) {
	s.wsHandlers[name] = handler
}

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
		s.log.Debug("websock subscribed to project",
			zap.String("userID", userID),
			zap.String("projectID", projectID))
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
		s.log.Debug("websock unsubscribed from project",
			zap.String("userID", userID),
			zap.String("projectID", projectID))
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
		return c.WriteJSON(
			map[string]interface{}{
				"type": "pong",
			},
		)
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
			reqLog := s.log.With(zap.String("userID", userID))
			reqLog.Debug("upgraded websocket connection")
			defer reqLog.Debug("closed websocket connection")
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
			stopped := make(chan error, 1)
			go handleWebSockMessages(
				ctx,
				c,
				messages,
				stopped,
				reqLog,
			)
			heartbeat := time.Now()
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case err := <-stopped:
					return errors.Wrap(err, "handleWebSockMessages")
				case <-time.After(pingInterval):
					if err := c.WriteJSON(
						map[string]interface{}{
							"type": "ping",
						},
					); err != nil {
						return err
					}
					if time.Since(heartbeat) > heartbeatTimeout {
						return errors.New("heartbeat timeout")
					}
				case msg := <-messages:
					obj := make(map[string]interface{})
					if err := json.Unmarshal(msg, &obj); err != nil {
						fmt.Printf("message='%s'\n", string(msg))
						return errors.Wrap(err, "unmarshal websocket message")
					}
					ty, ok := obj["type"].(string)
					if !ok {
						return errors.New("invalid message")
					} else if ty == "pong" {
						heartbeat = time.Now()
						continue
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

func handleWebSockMessages(
	ctx context.Context,
	c *websocket.Conn,
	messages chan<- []byte,
	err chan<- error,
	log *zap.Logger,
) {
	err <- func() error {
		defer close(messages)
		for {
			ty, buf, err := c.ReadMessage()
			if err != nil {
				return err
			}
			switch ty {
			case websocket.CloseMessage:
				_ = c.Close()
				return nil
			case websocket.PingMessage:
				if err := c.WriteMessage(
					websocket.PongMessage,
					nil,
				); err != nil {
					return err
				}
			case websocket.TextMessage:
				select {
				case <-ctx.Done():
					return ctx.Err()
				case messages <- buf:
					continue
				default:
					log.Warn("message dropped")
				}
			default:
				return fmt.Errorf("unexpected message type: %d", ty)
			}
		}
	}()
}
