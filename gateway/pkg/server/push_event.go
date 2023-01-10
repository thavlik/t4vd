package server

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/gateway/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) PushEvent(
	ctx context.Context,
	req api.Event,
) (*api.Void, error) {
	if len(req.ProjectIDs) == 0 {
		return nil, errors.New("missing projectIDs")
	}
	body, err := json.Marshal(&req)
	if err != nil {
		return nil, errors.Wrap(err, "json")
	}
	if err := s.pub.Publish(
		ctx,
		gatewayTopic,
		body,
	); err != nil {
		return nil, errors.Wrap(err, "publisher.Publish")
	}
	return &api.Void{}, nil
}

func (s *Server) getSubscriptions(projectIDs []string) (subs []*Subscription) {
	s.subsL.Lock()
	for sub := range s.subs {
		for _, projectID := range projectIDs {
			if sub.projectID == projectID {
				subs = append(subs, sub)
				break
			}
		}
	}
	s.subsL.Unlock()
	return
}

func (s *Server) pushEventLocal(req api.Event) error {
	subs := s.getSubscriptions(req.ProjectIDs)
	for _, sub := range subs {
		select {
		case sub.ch <- []byte(req.Payload):
			continue
		default:
			// TODO: should we close the subscription here?
			s.log.Warn("discarding event due to full stream",
				zap.String("userID", sub.userID))
		}
	}
	return nil
}
