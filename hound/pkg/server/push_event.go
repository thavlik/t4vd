package server

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	gateway "github.com/thavlik/t4vd/gateway/pkg/api"
)

func (s *Server) PushEvent(
	ctx context.Context,
	ty string,
	payload interface{},
	projectIDs []string,
) error {
	pl, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}
	body, err := json.Marshal(&EventWrapper{
		Type:    ty,
		Payload: string(pl),
	})
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}
	if _, err := s.gateway.PushEvent(
		context.Background(),
		gateway.Event{
			ProjectIDs: projectIDs,
			Payload:    string(body),
		},
	); err != nil {
		return errors.Wrap(err, "gateway.PushEvent")
	}
	return nil
}
