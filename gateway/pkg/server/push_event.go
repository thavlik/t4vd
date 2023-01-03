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
	if err := s.pub.Publish(body); err != nil {
		return nil, errors.Wrap(err, "pubsub.Publish")
	}
	return &api.Void{}, nil
}

func (s *Server) getSubscriptions(projectIDs []string) (subs []*Subscription) {
	s.subsL.Lock()
	for sub := range s.subs {
		found := false
		for _, projectID := range projectIDs {
			if sub.projectID == projectID {
				found = true
				break
			}
		}
		if found {
			subs = append(subs, sub)
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
			s.log.Warn("client event stream is full",
				zap.String("userID", sub.userID))
		}
	}
	return nil
}
